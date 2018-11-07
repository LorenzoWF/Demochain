package main

import (
	"flag"
	"log"
	"os"
	golog "github.com/ipfs/go-log"
	gologging "github.com/whyrusleeping/go-logging"
	"github.com/joho/godotenv"
	core "core"
	network "network"
	host "host"
)

func main() {
	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	//GRAVA LOG, VER AONDE DA PRA VER ESSA MERDA
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info


	peerIdParm := flag.Int("p", 0, "wait for incoming connections")
	flag.Parse()

	if *peerIdParm == 0 {
		log.Fatal("Please provide a id of peer to bind on with -p") //INTERROMPE A EXECUCAO
	}

	switch *peerIdParm {
		case 1:
			err := godotenv.Load("peer1.env")
			if err != nil {
				log.Fatal(err)
			}

		case 2:
			err := godotenv.Load("peer2.env")
			if err != nil {
				log.Fatal(err)
			}

		case 3:
			err := godotenv.Load("peer3.env")
			if err != nil {
				log.Fatal(err)
			}
	}

	//Carrega as configuracoes do .env
	nodeIP 								 := os.Getenv("IP")
	nodePort 							 := os.Getenv("TCP_PORT")
	nodeNetworkName 			 := os.Getenv("NETWORK_NAME")
	nodePathPrivateKey 		 := os.Getenv("PATH_PRIVATE_KEY")
	nodeCryptographicType  := os.Getenv("CRYPTOGRAPHIC_TYPE")
	nodeCryptographicBits  := os.Getenv("CRYPTOGRAPHIC_BITS")
	nodeEDNodeTarget 			 := os.Getenv("EDGE_NODE_TARGET")
	nodePathBlockchainFile := os.Getenv("PATH_BLOCKCHAIN_FILE")
	nodeHLNodes						 := os.Getenv("HL_NODES")
	nodeConsensus					 := os.Getenv("CONSENSUS")
	nodeDifficulty				 := os.Getenv("DIFFICULTY")

	node := core.NodeLoad(nodeIP, nodePort, nodeNetworkName, nodePathPrivateKey, nodeCryptographicType, nodeCryptographicBits, nodeEDNodeTarget, nodePathBlockchainFile, nodeHLNodes, nodeConsensus, nodeDifficulty)

	// Make a host that listens on the given multiaddress
	peer, err := network.MakeBasicHost(node)
	if err != nil {
		log.Fatal(err)
	}

	fullAddr := network.MakeFullAddr(peer)
	log.Println(fullAddr)

	host.Init(node, peer)

	host.BlockchainLoad()

	if *peerIdParm == 1 {
			host.HostHandler() //ESCUTA OUTRAS CONEXOES
	} else {
			host.Connect() //SOMENTE CONECTA
	}

	select {}

}
