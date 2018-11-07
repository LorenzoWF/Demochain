package storage

import (
	"bufio"
	"log"
	"encoding/json"
	"fmt"
	"time"
	"sync"
	"os"
	//"bytes"
	"strconv"
	"strings"
	"io/ioutil"
	"github.com/davecgh/go-spew/spew"
	"context"
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
//var RW *bufio.ReadWriter
var SendStatus bool

func Init(node core.Node, peer host.Host) {
		Node = node
		Host = peer
}

func HostHandler() {
	SendStatus = false
	Host.SetStreamHandler(Node.NetworkName, HandleStream)
}

func HandleStream(s net.Stream) {

	//log.Println(s.Conn().RemotePeer().MatchesPublicKey())

	var pertence bool
	var permiss int
	log.Println(permiss)
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

		go ReadData(rw)
		go WriteData(rw)
		SendStatus = true

		// stream 's' will stay open until you close it (or the other side closes it).
	} else {
		fmt.Printf("\x1b[36m%s\x1b[0m", "Peer Rejeitado.\n")
	}
}

/*func sendAllBlocks(rw *bufio.ReadWriter) {
	//for {
		//JA ESTAO MINERADOS
		for _, bloco := range Blockchain {
				bytesBlocoMinerado, _ := json.Marshal(bloco)
				sendBlock(rw, bytesBlocoMinerado)
		}
	//}
}*/

func BlockchainLoad() {
	bcFile, err := ioutil.ReadFile(Node.PathBlockchainFile)
	err = json.Unmarshal(bcFile, &Blockchain)
	if err != nil {
		//log.Info(err)
		log.Println("BC não carregada")
	}

	if len(Blockchain) == 0 {
		Blockchain = append(Blockchain, core.GenerateGenesisBlock(Node.Consensus, Node.Difficulty))
		bytesFile, _ := json.Marshal(Blockchain)
		writeFile(bytesFile)
	}
}

/*func Disconnect() {
	SendStatus = false
	Host.Close()
}*/

func Connect() {
	//HostHandler() //NAO E necessario

	SendStatus = false

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

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go ReadData(rw)
	go WriteData(rw)
	//RW = rw
	SendStatus = true
}

/*func sendHandler(rw *bufio.ReadWriter) {

	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)
		sendData = strings.Replace(sendData, "\r", "", -1)
		bpm, _ := strconv.Atoi(sendData)

		Miner(rw, bpm)
	}
}

func Miner(rw *bufio.ReadWriter, data int) {
		newBlock := core.GenerateBlock(Blockchain[len(Blockchain)-1], data, Node.Consensus, Node.Difficulty)

		if core.IsBlockValid(newBlock, Blockchain[len(Blockchain)-1], Node.Consensus) {
			mutex.Lock()
			Blockchain = append(Blockchain, newBlock)
			bytesFile, err := json.Marshal(Blockchain)
			writeFile(bytesFile)
			mutex.Unlock()

			bytes, err := json.Marshal(newBlock)
			if err != nil {
				log.Println(err)
			}

			spew.Dump(newBlock)

			//if SendStatus == true {
					go sendBlock(rw, bytes)
			//}

		} else {
			log.Println("Bloco invalido")
		}
}*/

/*func receiveHandler(rw *bufio.ReadWriter) {
	//LOOP PARA VERIFICAR SE HOUVE MINERACAO EM OUTRO NODOS
	for {
			str, err := rw.ReadString('\n')
			if err != nil {
				fmt.Printf("\x1b[33m%s\x1b[0m> ", "Error! Peer Diconencted.\n")
			}

			if str == "" {
				return
			}

			if str != "\n" {

				var blockRecebido core.Block
				var blockAtual core.Block
				err := json.Unmarshal([]byte(str), &blockRecebido);

				if err != nil {
					log.Fatal(err)
				}

				mutex.Lock()

				if len(Blockchain) > 0 {
				 	blockAtual = Blockchain[len(Blockchain)-1]
				} else {
					blockAtual.Index = 0
					blockAtual.Hash = ""
				}

				addBlock := false

				if blockRecebido.Index == 0 && blockAtual.Index == 0 {
					addBlock = true
				} else {
					if blockRecebido.Index > blockAtual.Index {
						if core.IsBlockValid(blockRecebido, blockAtual, Node.Consensus) { //VERIFICACAO DE HASH
							addBlock = true
						}
					}
				}

				if addBlock == true {
					Blockchain = append(Blockchain, blockRecebido)
					bytesFile, err := json.Marshal(Blockchain)
					bytes, err := json.MarshalIndent(blockRecebido, "", "  ")
					if err != nil {
						log.Fatal(err)
					}
					writeFile(bytesFile)
					//sendBlock(RW, bytes)
					// Green console color: 	\x1b[32m
					// Reset console color: 	\x1b[0m
					fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
				}

				mutex.Unlock()
			}
	}
}*/

/*func sendBlock(rw *bufio.ReadWriter, bytes []byte) {
	mutex.Lock()
	rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
	rw.Flush()
	mutex.Unlock()
}*/


func writeFile(bytes []byte) {
	err := ioutil.WriteFile(Node.PathBlockchainFile, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
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

			if err != nil {
				log.Fatal(err)
			}

			mutex.Lock()
			if len(chain) > len(Blockchain) {
				if core.IsBlockValid(chain[len(chain)-1], Blockchain[len(Blockchain)-1], Node.Consensus) { //VERIFICACAO DE HASH
					Blockchain = chain
					bytes, err := json.MarshalIndent(Blockchain, "", "  ")
					bytesArquivo, err := json.Marshal(Blockchain)
					if err != nil {
						log.Fatal(err)
					}
					writeFile(bytesArquivo)
					// Green console color: 	\x1b[32m
					// Reset console color: 	\x1b[0m
					fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
				}
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

	//TROCAR POR ARQUIVO DEPOIS
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)
		sendData = strings.Replace(sendData, "\r", "", -1)
		data, err := strconv.Atoi(sendData)
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m", "Error! Invalid Data.\n")
			continue
		}
		//newHash, nonce := core.MinerBlock(data, Node.Consensus, Node.Difficulty)
		var newHash, prevHash, nonce string
		var done bool
		if Node.Consensus == 1 {
			for i := 0; ; i++ {
					prevHash = Blockchain[len(Blockchain)-1].Hash
					done, newHash, nonce = core.MinerBlockLoop(i, data, prevHash, Node.Consensus, Node.Difficulty)
					if done == false {
						if prevHash != Blockchain[len(Blockchain)-1].Hash {
								i = 0 //SE ATUALIZOU O BLOCO, O LOOP FAZ DE NOVO
						}

						continue
					} else {
						break
					}
			}
		}
		//BlockchainLoad() //ATUALIZA BLOCKCHAIN

		newBlock := core.GenerateBlock(Blockchain[len(Blockchain)-1], data, newHash, nonce, Node.Consensus, Node.Difficulty)

		if core.IsBlockValid(newBlock, Blockchain[len(Blockchain)-1], Node.Consensus) {
			mutex.Lock()
			Blockchain = append(Blockchain, newBlock)
			mutex.Unlock()
		} else {
			log.Println("Bloco invalido")
			//log.Println("Hash antigo: %s", hashAntigo)
			//log.Println("Hash novo: %s", hashNovo)
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
