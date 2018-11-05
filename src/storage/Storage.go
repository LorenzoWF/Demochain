package storage

import (
	"bufio"
	"log"
	"encoding/json"
	"fmt"
	"time"
	"sync"
	"os"
	"strconv"
	"strings"
	"io/ioutil"
	"context"
	"github.com/davecgh/go-spew/spew"
	net "github.com/libp2p/go-libp2p-net"
	core "core"
	host "github.com/libp2p/go-libp2p-host"
	ma "github.com/multiformats/go-multiaddr"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

var Blockchain []core.Block
var Node core.Node
var Host host.Host
var mutex = &sync.Mutex{}

func HandleStream(s net.Stream) {

	//log.Println(s.Conn().RemotePeer().MatchesPublicKey())

	var pertence bool
	var permiss int
	HLNodes := Node.HLNodes

	pertence = false

	for _, HLNode := range HLNodes {
			if s.Conn().RemotePeer().Pretty() == HLNode.HLTarget {
					pertence = true
					permiss  = HLNode.Permiss
					break
			}
	}

	if pertence == true {
		fmt.Printf("\x1b[36m%s\x1b[0m", "Peer Conencted.\n")

		log.Println("REDE: ", s.Protocol())
		log.Println("Eu sou: ", s.Conn().LocalPeer().Pretty())
		log.Println("Recebendo conexao de: ", s.Conn().RemotePeer().Pretty())

		// Create a buffer stream for non blocking read and write.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		if permiss == 1 { //NAO VOU LER NADA DO QUE ELE ESCREVER
			go ReadData(rw)
		}
		go WriteData(rw)

		// stream 's' will stay open until you close it (or the other side closes it).
	} else {
		fmt.Printf("\x1b[36m%s\x1b[0m", "Peer Rejeitado.\n")
	}
}

func TesteNotificacao(t net.Stream) {
	log.Println("ENTROU NA TESTE")
}

func ReadData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Printf("\x1b[33m%s\x1b[0m> ", "Error! Peer Diconencted.\n")
		}

		if str == "" {
			return
		}
		if str != "\n" {
			chain := make([]core.Block, 0)
			err := json.Unmarshal([]byte(str), &chain);

			//log.Println(len(chain))

			if err != nil {
				log.Fatal(err)
			}

			mutex.Lock()
			if len(chain) > len(Blockchain) {
				Blockchain = chain
				bytes, err := json.MarshalIndent(Blockchain, "", "  ")
				if err != nil {
					log.Fatal(err)
				}
				writeFile(bytes)
				// Green console color: 	\x1b[32m
				// Reset console color: 	\x1b[0m
				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
			}

			mutex.Unlock()
		}
	}
}

func WriteData(rw *bufio.ReadWriter) {
	//IMPRIMI OS DADOS QUE FOI INSERIDO NO MESMO TERMINAL
	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Println(err)
			}
			mutex.Unlock()

			mutex.Lock()
			writeFile(bytes)
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			mutex.Unlock()
		}
	}()

	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)
		sendData = strings.Replace(sendData, "\r", "", -1)

		if sendData == "id" {
			log.Println(Host.ID())
			continue
		}

		if sendData == "addrs" {
			log.Println(Host.Addrs())
			continue
		}

		if sendData == "node" {
			spew.Dump(Node)
			continue
		}

		if sendData == "connect" {
			connect()
			continue
		}

		if sendData == "disconnect" {
			disconnect()
			continue
		}


		bpm, err := strconv.Atoi(sendData)
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m", "Error! Invalid Data.\n")
			continue
		}
		newBlock := core.GenerateBlock(Blockchain[len(Blockchain)-1], bpm)

		if core.IsBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
			mutex.Lock()
			Blockchain = append(Blockchain, newBlock)
			mutex.Unlock()
		}

		bytes, err := json.Marshal(Blockchain)
		if err != nil {
			log.Println(err)
		}

		spew.Dump(Blockchain)

		mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		mutex.Unlock()
	}
}

func writeFile(bytes []byte) {
	err := ioutil.WriteFile(Node.PathBlockchainFile, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func BlockchainLoad() {
	bcFile, err := ioutil.ReadFile(Node.PathBlockchainFile)
	err = json.Unmarshal(bcFile, &Blockchain)
	if err != nil {
		//log.Info(err)
		log.Println("BC não carregada")
	}

	if len(Blockchain) == 0 {
		generateGenesisBlock()
	}
}

func generateGenesisBlock() {
	t := time.Now()
	genesisBlock := core.Block{}
	genesisBlock = core.Block{0, t.String(), 0, core.CalculateHash(genesisBlock), "", ""}
	Blockchain = append(Blockchain, genesisBlock)
}

func HostLoad() {
	if Node.ELTarget == "" {
		// Set a stream handler on host A. /p2p/1.0.0 is
		log.Println("listening for connections")
		// a user-defined protocol name.

		Host.SetStreamHandler(Node.NetworkName, HandleStream)

		select {} // hang forever
		//This is where the listener code ends
	} else {
		connect()
	}
}

func disconnect() {
	Host.Close()
}

func connect() {
	Host.SetStreamHandler(Node.NetworkName, HandleStream)

	// The following code extracts target's peer ID from the
	ipfsaddr, err := ma.NewMultiaddr(Node.ELTarget)
	// given multiaddress
	if err != nil {
		log.Fatalln(err)
	}

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		log.Fatalln(err)
	}

	peerid, err := peer.IDB58Decode(pid)
	if err != nil {
		log.Fatalln(err)
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	Host.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

	log.Println("opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /p2p/1.0.0 protocol

	s, err := Host.NewStream(context.Background(), peerid, Node.NetworkName)
	if err != nil {
		log.Fatalln(err)
	}

	// Create a buffered stream so that read and writes are non blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// Create a thread to read and write data.
	//go WriteData(rw) //ESCREVE
	go ReadData(rw) //ESCUTA

	select {} // hang forever
}
