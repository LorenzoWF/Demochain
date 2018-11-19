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

type DemoHost struct {
		Blockchain []core.Block
		Node core.Node
		Host host.Host
		SendStatus bool
		Conns []Conn
		Wg sync.WaitGroup
		Wg2 sync.WaitGroup
		Controle2 bool
		ConnectStatus bool
		Mutex *sync.Mutex
		//Mensagem string
		Mensagem MsgStruct
		Retornos []ReturnsStruct
}

type MsgStruct struct {
		Target string
		TypeM int
		Content string
		BlockchainMsg []core.Block
}

type ReturnsStruct struct {
		Target string
		Content string
}

func DemoHostCreate(node core.Node, host host.Host) (*DemoHost) {
		demoHost := DemoHost{}
		demoHost.Node = node
		demoHost.Host = host
		demoHost.SendStatus = false
		demoHost.ConnectStatus = false
		demoHost.Mutex = &sync.Mutex{}
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
			conn := Conn{s.Conn().RemotePeer().Pretty(), 0, c}
			demoHost.Conns = append(demoHost.Conns, conn)
			demoHost.SendStatus = true
			log.Println("Config send")
			go demoHost.sendHandler(rw, c)
		}

		//PRIMEIRO ENVIO SEM TUNELAMENTO
		demoHost.prepareMsg(demoHost.Node.GetTarget(), "blockchain", 0, true)
		bytes, err := json.Marshal(demoHost.Mensagem)
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		}

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

	if len(demoHost.Node.GetHLNodes()) > 0 || demoHost.Node.GetPermissioned() == 0 {
			demoHost.hostHandler()
	}

	if demoHost.Node.GetELTarget() != "" {
		ipfsaddr, err := ma.NewMultiaddr(demoHost.Node.GetELTarget())
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		}

		targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		demoHost.Host.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")

		s, err := demoHost.Host.NewStream(context.Background(), peerid, demoHost.Node.GetNetworkName())
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		}

		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		go demoHost.receiveHandler(rw)
		c := make(chan int, 2)

		subs := strings.Split(demoHost.Node.GetELTarget(), "/")
		conn := Conn{subs[6], 0, c}
		demoHost.Conns = append(demoHost.Conns, conn)

		demoHost.SendStatus = true
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
		if demoHost.Controle2 == true {
				demoHost.Wg2.Done()
				demoHost.Controle2 = false
		}
		log.Println("waiting...")
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Printf("\x1b[33m%s\x1b[0m> ", "Error! Peer Diconencted.\n")
		}

		if str == "" {
			return
		}

		str = strings.Replace(str, "\n", "", -1)
		str = strings.Replace(str, "\r", "", -1)

		var mensagemRecebida MsgStruct
		err = json.Unmarshal([]byte(str), &mensagemRecebida);

		if err != nil {
			log.Fatal(err)
		}

		if mensagemRecebida.Content == "verifyPBFT" {
			for _, conn := range demoHost.Conns {
					if conn.Target == mensagemRecebida.Target && conn.Target != "" {
						demoHost.prepareMsg(demoHost.Node.GetTarget(), "returnPBFT", 0, false)
						demoHost.Wg.Add(1)
						conn.Canal <- 0
						demoHost.Wg.Wait()
					}
			}

			continue
		}

		if mensagemRecebida.Content == "returnPBFT" {
			for _, conn := range demoHost.Conns {
					if conn.Target == mensagemRecebida.Target && conn.Target != "" {

						//IMPLEMENTAR VALIDACAO
						var retorno ReturnsStruct
						retorno = ReturnsStruct{conn.Target, "validateSuccess"}
						demoHost.Retornos = append(demoHost.Retornos, retorno)

						demoHost.Wg.Done()
					}
			}
			continue
		}

		if mensagemRecebida.Content != "\n" {
			demoHost.Wg2.Add(1)
			demoHost.Controle2 = true
			demoHost.Mutex.Lock()
			//chain := make([]core.Block, 0)
			//err := json.Unmarshal([]byte(str), &chain);
			chain := mensagemRecebida.BlockchainMsg

			if err != nil {
				log.Fatal(err)
			}

			if len(chain) > len(demoHost.Blockchain) {
				for len(demoHost.Blockchain) < len(chain) {
					if chain[len(demoHost.Blockchain)].IsBlockValid(demoHost.Blockchain[len(demoHost.Blockchain)-1]) { //VERIFICACAO DE HASH
						//demoHost.Blockchain = chain //NAO TEM PROBLEMA EM SUBSTITUIR O GENESIS
						demoHost.Blockchain = append(demoHost.Blockchain, chain[len(demoHost.Blockchain)]) //ADD UM POR 1
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

						for _, conn := range demoHost.Conns {
								if conn.Target != demoHost.Blockchain[len(demoHost.Blockchain)-1].GetTarget() && conn.Target != "" {
									log.Println("send to ", conn.Target)

									demoHost.prepareMsg(demoHost.Node.GetTarget(), "blockchain", 0, true)

									demoHost.Wg.Add(1)
									conn.Canal <- demoHost.Blockchain[len(demoHost.Blockchain)-1].GetIndex()
									demoHost.Wg.Wait()
								}
						}
					} else {
						break
					}
				}
			}
			demoHost.Mutex.Unlock()
		}
	}
}

func (demoHost *DemoHost) sendHandler(rw *bufio.ReadWriter, c chan int) {
	for {
		<-c //AGUARDA RECEBER AUTORIZACAO

		bytes, err := json.Marshal(demoHost.Mensagem)

		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		}

		/***ENVIA***/
		//mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		//mutex.Unlock()

		if demoHost.Mensagem.TypeM == 0 { //0 - Assync - 1 - Sync
				demoHost.Wg.Done()
		}
	}
}


func (demoHost *DemoHost) ProcessBlock(data string) {
	//time.Sleep(2 * time.Second)
	for {
		var newHash, prevHash, nonce string
		var done, validateSuccess bool

		validateSuccess = false

		if demoHost.Node.GetConsensus().GetType() == 1 { //Proof of Work
			for i := 0; ; i++ {
					//TESTAR COM MUTEXT, QUEM SABE TIRAR O PREV HASH
					prevHash = demoHost.Blockchain[len(demoHost.Blockchain)-1].GetHash()
					done, newHash, nonce = core.MakeBlockLoop(i, data, prevHash, demoHost.Node.GetConsensus())
					if done == false {
						if prevHash != demoHost.Blockchain[len(demoHost.Blockchain)-1].GetHash() {
								i = 0
						}
						continue
					} else {
						validateSuccess = true
						break
					}
			}
		} else if demoHost.Node.GetConsensus().GetType() ==  2 { //Proof of Stake
				prevHash = demoHost.Blockchain[len(demoHost.Blockchain)-1].GetHash()
				done, newHash = core.ProcessHash(data, prevHash, demoHost.Node.GetConsensus())
				demoHost.Retornos = demoHost.Retornos[:0]
				demoHost.prepareMsg(demoHost.Node.GetTarget(), "verifyPOS", 1, false)
				demoHost.sendMsg(0)
				for _, retorno := range demoHost.Retornos {
						if retorno.Content == "validateSuccess" {
								validateSuccess = true
								break;
						}
				}
		} else if demoHost.Node.GetConsensus().GetType() ==  3 { //PBFT
				prevHash = demoHost.Blockchain[len(demoHost.Blockchain)-1].GetHash()
				done, newHash = core.ProcessHash(data, prevHash, demoHost.Node.GetConsensus())
				demoHost.Retornos = demoHost.Retornos[:0]
				demoHost.prepareMsg(demoHost.Node.GetTarget(), "verifyPBFT", 1, false)
				demoHost.sendMsg(0)
				numberAccept := 0
				numberNoAccept := 0
				for _, retorno := range demoHost.Retornos {
						if retorno.Content == "validateSuccess" {
								numberAccept += 1
						} else {
								numberNoAccept += 1
						}
				}
				validateSuccess = core.ValidatePBFT(numberAccept, numberNoAccept)
		} else {
			log.Println("Consensus invalid")
			return
		}

		if validateSuccess == false {
				return
		}

		log.Println("Refresh Blocks")
		demoHost.Wg2.Wait()

		newBlock := core.GenerateBlock(demoHost.Blockchain[len(demoHost.Blockchain)-1], data, newHash, nonce, demoHost.Node.GetConsensus(), demoHost.Node.GetTarget())

		if newBlock.IsBlockValid(demoHost.Blockchain[len(demoHost.Blockchain)-1]) {
			demoHost.Mutex.Lock()
			demoHost.Blockchain = append(demoHost.Blockchain, newBlock)

			bytes, err := json.Marshal(demoHost.Blockchain)
			if err != nil {
				//log.Println(err)
				fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
			}

			demoHost.writeFile(bytes)

			spew.Dump(newBlock)

			demoHost.prepareMsg(demoHost.Node.GetTarget(), "blockchain", 0, true)
			demoHost.sendMsg(newBlock.GetIndex())

			demoHost.Mutex.Unlock()
		} else {
			log.Println("Invalid block!")
		}
		return //NECESSARIO
	}
}

func (demoHost *DemoHost) sendMsg(canal int) {
	if demoHost.SendStatus == true {
			for _, conn := range demoHost.Conns {
					if conn.Target != "" {
							demoHost.Wg.Add(1)
							conn.Canal <- conn.Index
							demoHost.Wg.Wait()
					}
			}
	}
}

func (demoHost *DemoHost) prepareMsg(target string, content string, typeM int, refreshBlockchainMsg bool)  {

		demoHost.Mensagem.Target = target
		demoHost.Mensagem.Content = content
		demoHost.Mensagem.TypeM = typeM

		if refreshBlockchainMsg == true {
				demoHost.Mensagem.BlockchainMsg = demoHost.Blockchain
		} else {
				demoHost.Mensagem.BlockchainMsg = demoHost.Mensagem.BlockchainMsg[:0]
		}
}
