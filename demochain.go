package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"time"
	"os"
	"strconv"

	golog "github.com/ipfs/go-log"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	gologging "github.com/whyrusleeping/go-logging"
	"github.com/joho/godotenv"


	protocol "github.com/libp2p/go-libp2p-protocol"

	core "core"
	network "network"
	storage "storage"
)

func main() {
	//LE AS VAR DO .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//Carrega as configuracoes
	nodeIP := os.Getenv("IP")

	var networkName protocol.ID
	networkName = protocol.ID(os.Getenv("NETWORK_NAME"))

	cryptographicType, _ := strconv.Atoi(os.Getenv("CRYPTOGRAPHIC_TYPE"))
	cryptographicBits, _ := strconv.Atoi(os.Getenv("CRYPTOGRAPHIC_BITS"))

	Node := core.Node{nodeIP, cryptographicType, cryptographicBits}

	t := time.Now() //PEGA A HORA ATUAL
	genesisBlock := core.Block{} //CRIA O GENESIS BLOCK, TIPO
	genesisBlock = core.Block{0, t.String(), 0, core.CalculateHash(genesisBlock), "", ""} //CRIA O GENESIS BLOCK

	storage.Blockchain = append(storage.Blockchain, genesisBlock) //ADICIONA NA BLOCKCHAIN

	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	//GRAVA LOG, VER AONDE DA PRA VER ESSA MERDA
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	// Parse options from the command line
	//ipNode := flag.String("i", "", "IP for connection")
	portNode := flag.Int("p", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	//seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *portNode == 0 {
		log.Fatal("Please provide a port to bind on with -p") //INTERROMPE A EXECUCAO
	}

	// Make a host that listens on the given multiaddress
	ha, err := network.MakeBasicHost(Node.IP, *portNode, []byte("TESTE"), Node.CryptographicType, Node.CryptographicBits)
	if err != nil {
		log.Fatal(err) //INTERROMPE A EXECUCAO
	}

	if *target == "" {
		// Set a stream handler on host A. /p2p/1.0.0 is
		log.Println("listening for connections")
		// a user-defined protocol name.

		ha.SetStreamHandler(networkName, storage.HandleStream)

		select {} // hang forever
		/**** This is where the listener code ends ****/
	} else {
		//ACHO QUE NAO PRECISA
		ha.SetStreamHandler(networkName, storage.HandleStream)

		// The following code extracts target's peer ID from the
		ipfsaddr, err := ma.NewMultiaddr(*target)
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
		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := ha.NewStream(context.Background(), peerid, networkName)
		if err != nil {
			log.Fatalln(err)
		}
		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go storage.WriteData(rw)
		go storage.ReadData(rw)

		//CRIAR PARTE QUE

		select {} // hang forever
	}
}
