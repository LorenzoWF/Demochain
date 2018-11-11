package network

import (
	"bufio"
	"log"
	"encoding/json"
	"fmt"
	//"time"
	"sync"
	//"os"
	//"bytes"
	//"strconv"
	"strings"
	"io/ioutil"
	"github.com/davecgh/go-spew/spew"
	"context"
	net "github.com/libp2p/go-libp2p-net"
	core "demochain/core"
	host "github.com/libp2p/go-libp2p-host"
	ma "github.com/multiformats/go-multiaddr"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

type conn struct {
		target string
		index int
		canal chan int
}

type DemoHost struct {
	Blockchain []core.Block
	Node core.Node
	Host host.Host

	sendStatus bool
	conns []conn

	wg sync.WaitGroup

	wg2 sync.WaitGroup
	controle2 bool
	ConnectStatus bool
	mutex *sync.Mutex
	//var canal2 chan int
}

func DemoHostCreate(node core.Node, host host.Host) (*DemoHost) {
		demoHost := DemoHost{}
		demoHost.Node = node
		demoHost.Host = host
		demoHost.sendStatus = false
		demoHost.ConnectStatus = false
		demoHost.mutex = &sync.Mutex{}
		demoHost.Blockchain = demoHost.Blockchain[:0]

		return &demoHost //return pointer
}

func (demoHost *DemoHost) GetBlockchain() ([]core.Block) {
		return demoHost.Blockchain
}

func (demoHost *DemoHost) GetNode() (core.Node) {
		return demoHost.Node
}

func (demoHost *DemoHost) hostHandler() {
	demoHost.Host.SetStreamHandler(demoHost.Node.GetNetworkName(), demoHost.handleStream)
}

func (demoHost *DemoHost) handleStream(s net.Stream) {
	var pertence bool
	pertence = false

	var permiss int

	HLNodes := demoHost.Node.GetHLNodes()

	if demoHost.Node.GetPermissioned() == 1 { //REDE PERMISSIONADA
		for _, HLNode := range HLNodes {
				if s.Conn().RemotePeer().Pretty() == HLNode.GetHLTarget() {
						pertence = true
						permiss  = HLNode.GetPermiss()
						break
				}
		}
		log.Println("Permiss = ", permiss)
	}

	if pertence == true || demoHost.Node.GetPermissioned() == 0 {
		fmt.Printf("\x1b[36m%s\x1b[0m", "Peer Conencted.\n")

		log.Println("Network Name: ", s.Protocol())
		log.Println("I am: ", s.Conn().LocalPeer().Pretty())
		log.Println("Receiving connection from: ", s.Conn().RemotePeer().Pretty())

		// Create a buffer stream for non blocking read and write.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		log.Println("Config receive")
		go demoHost.receiveHandler(rw)

		if permiss == 1 || demoHost.Node.GetPermissioned() == 0 { //PERMISSAO PARA ESCRITA
			c := make(chan int, 2)
			conn := conn{s.Conn().RemotePeer().Pretty(), 0, c}
			demoHost.conns = append(demoHost.conns, conn)
			demoHost.sendStatus = true
			log.Println("Config send")
			go demoHost.sendHandler(rw, c)
		}

		//PRIMEIRO ENVIO SEM TUNELAMENTO
		bytes, err := json.Marshal(demoHost.Blockchain)
		if err != nil {
			//log.Println(err)
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		}

		/***ENVIA***/
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()

		demoHost.ConnectStatus = true

		// stream 's' will stay open until you close it (or the other side closes it).
	} else {
		fmt.Printf("\x1b[36m%s\x1b[0m", "Peer Rejeitado.\n")
	}
}

func (demoHost *DemoHost) BlockchainClean() {
	demoHost.Blockchain = demoHost.Blockchain[:0]
}

func (demoHost *DemoHost) BlockchainLoad() {
//func BlockchainLoad(demoHost *DemoHost) {
	bcFile, err := ioutil.ReadFile(demoHost.Node.GetPathBlockchainFile())
	err = json.Unmarshal(bcFile, &demoHost.Blockchain)
	if err != nil {
		//log.Info(err)
		//log.Println("Blockchain not loaded")
		fmt.Printf("\x1b[33m%s\x1b[0m\n", "File is empty")
	}

	if len(demoHost.Blockchain) == 0 {
		demoHost.Blockchain = append(demoHost.Blockchain, core.GenerateGenesisBlock(demoHost.Node.GetConsensus(), demoHost.Node.GetTarget()))
		bytesFile, err := json.Marshal(demoHost.Blockchain)

		if err != nil {
			log.Println(err)
		}

		demoHost.writeFile(bytesFile)
	}
}

//IMPLEMENTAR
func (demoHost *DemoHost) BlockchainRequest() {

}

func (demoHost *DemoHost) GetBlock(indexBlock int) (core.Block) {
	var block core.Block
	if len(demoHost.Blockchain) >= (indexBlock + 1) {
		block = demoHost.Blockchain[indexBlock]
	}
	return block
}

//FAZER FUNCAO DISCONNECT
func (demoHost *DemoHost) Disconnect() {
	demoHost.Host.Close()
	demoHost.ConnectStatus = false
}

func (demoHost *DemoHost) Connect() {

	if len(demoHost.Node.GetHLNodes()) > 0 {
			demoHost.hostHandler() //NAO E necessario
	}

	if demoHost.Node.GetELTarget() != "" { //INSERIR TRATAMENTO
		// The following code extracts target's peer ID from the
		ipfsaddr, err := ma.NewMultiaddr(demoHost.Node.GetELTarget())
		// given multiaddress
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
			//log.Fatalln(err)
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			//log.Fatalln(err)
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			//log.Fatalln(err)
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		}

		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it
		demoHost.Host.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := demoHost.Host.NewStream(context.Background(), peerid, demoHost.Node.GetNetworkName())
		if err != nil {
			//log.Fatalln(err)
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		}

		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		go demoHost.receiveHandler(rw)
		c := make(chan int, 2)

		subs := strings.Split(demoHost.Node.GetELTarget(), "/")
		conn := conn{subs[6], 0, c}
		demoHost.conns = append(demoHost.conns, conn)

		demoHost.sendStatus = true
		demoHost.ConnectStatus = true

		go demoHost.sendHandler(rw, c)
	}
}

func (demoHost *DemoHost) writeFile(bytes []byte) {
	//log.Println(string(bytes))
	err := ioutil.WriteFile(demoHost.Node.GetPathBlockchainFile(), bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}


func (demoHost *DemoHost) receiveHandler(rw *bufio.ReadWriter) {
	for {
		if demoHost.controle2 == true {
				demoHost.wg2.Done()
				demoHost.controle2 = false
		}
		log.Println("waiting...")
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Printf("\x1b[33m%s\x1b[0m> ", "Error! Peer Diconencted.\n")
		}

		if str == "" {
			return
		}
		if str != "\n" {
			demoHost.wg2.Add(1)
			demoHost.controle2 = true
			demoHost.mutex.Lock()
			chain := make([]core.Block, 0)
			err := json.Unmarshal([]byte(str), &chain);

			if err != nil {
				log.Fatal(err)
			}

			if len(chain) > len(demoHost.Blockchain) {
				for len(demoHost.Blockchain) < len(chain) {
					if chain[len(demoHost.Blockchain)].IsBlockValid(demoHost.Blockchain[len(demoHost.Blockchain)-1]) { //VERIFICACAO DE HASH
						demoHost.Blockchain = chain //NAO TEM PROBLEMA EM SUBSTITUIR O GENESIS
						//bytes, err := json.MarshalIndent(Blockchain, "", "  ")
						bytes, err := json.MarshalIndent(demoHost.Blockchain[len(demoHost.Blockchain)-1], "", "  ")
						bytesArquivo, err := json.Marshal(demoHost.Blockchain)
						if err != nil {
							log.Fatal(err)
						}
						demoHost.writeFile(bytesArquivo)
						// Green console color: 	\x1b[32m
						// Reset console color: 	\x1b[0m
						fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))

						for _, conn := range demoHost.conns {
								if conn.target != demoHost.Blockchain[len(demoHost.Blockchain)-1].GetTarget() && conn.target != "" {
									log.Println("send to ", conn.target)
									demoHost.wg.Add(1)
									conn.canal <- demoHost.Blockchain[len(demoHost.Blockchain)-1].GetIndex()
									demoHost.wg.Wait()
								}
						}
					} else {
						break
					}
				}
			}

			demoHost.mutex.Unlock()
		}
	}
}

func (demoHost *DemoHost) sendHandler(rw *bufio.ReadWriter, c chan int) {
	for {
		v := <-c //AGUARDA RECEBER AUTORIZACAO
		log.Println("Sending block ", v)

		bytes, err := json.Marshal(demoHost.Blockchain)
		if err != nil {
			//log.Println(err)
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		}

		/***ENVIA***/
		//mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		//mutex.Unlock()
		log.Println("Block submitted")
		demoHost.wg.Done()
	}
}


func (demoHost *DemoHost) ProcessBlock(data string) {
	//time.Sleep(2 * time.Second)
	for {
		var newHash, prevHash, nonce string
		var done bool

		if demoHost.Node.GetConsensus().GetType() == 1 {
			for i := 0; ; i++ {
					prevHash = demoHost.Blockchain[len(demoHost.Blockchain)-1].GetHash()
					done, newHash, nonce = core.MinerBlockLoop(i, data, prevHash, demoHost.Node.GetConsensus())
					if done == false {
						if prevHash != demoHost.Blockchain[len(demoHost.Blockchain)-1].GetHash() {
								i = 0 //SE ATUALIZOU O BLOCO, O LOOP FAZ DE NOVO
						}
						continue
					} else {
						break
					}
			}
		}

		log.Println("Refresh Blocks")
		demoHost.wg2.Wait()

		newBlock := core.GenerateBlock(demoHost.Blockchain[len(demoHost.Blockchain)-1], data, newHash, nonce, demoHost.Node.GetConsensus(), demoHost.Node.GetTarget())

		if newBlock.IsBlockValid(demoHost.Blockchain[len(demoHost.Blockchain)-1]) {
			demoHost.mutex.Lock()
			demoHost.Blockchain = append(demoHost.Blockchain, newBlock)

			bytes, err := json.Marshal(demoHost.Blockchain)
			if err != nil {
				//log.Println(err)
				fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
			}

			demoHost.writeFile(bytes)

			spew.Dump(newBlock)

			if demoHost.sendStatus == true {
					for _, conn := range demoHost.conns {
							if conn.target != "" {
									demoHost.wg.Add(1)
									conn.canal <- newBlock.GetIndex()
									demoHost.wg.Wait()
							}
					}
			}
			demoHost.mutex.Unlock()
		} else {
			log.Println("Invalid block!")
		}
		return //NECESSARIO
	}
}
