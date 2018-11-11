package main

import (
	"flag"
	"log"
	"os"
	golog "github.com/ipfs/go-log"
	gologging "github.com/whyrusleeping/go-logging"
	"github.com/joho/godotenv"

	//"time"

	"bufio"
	"strings"
	//"strconv"
	"fmt"
	//"encoding/json"

	core "demochain/core"
	network "demochain/network"
)

func main() {
	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	dir, _ := os.Getwd()
	sep := string(os.PathSeparator)
	//log.Println(dir)
	//log.Println(string(os.PathSeparator))

	peerIdParm := flag.String("c", "", "wait for incoming connections")
	flag.Parse()

	if *peerIdParm == "" {
		log.Fatal("Please provide a path of peer config to bind on with -c") //INTERROMPE A EXECUCAO
	}

	path := dir + sep + *peerIdParm + sep + *peerIdParm + ".env"

	err := godotenv.Load(path)
	if err != nil {
		log.Fatal(err)
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

	nodePathPrivateKey = dir + sep + *peerIdParm + sep + nodePathPrivateKey
	nodePathBlockchainFile = dir + sep + *peerIdParm + sep + nodePathBlockchainFile

	node := core.NodeLoad(nodeIP, nodePort, nodeNetworkName, nodePathPrivateKey, nodeCryptographicType, nodeCryptographicBits, nodeEDNodeTarget, nodePathBlockchainFile, nodeHLNodes, nodeConsensus, nodeDifficulty)

	host, err := network.MakeBasicHost(&node)
	log.Println(node.GetTarget())

	if err != nil {
		log.Fatal(err)
		log.Println(host)
	}

	fullAddr := network.MakeFullAddr(host)
	log.Println(fullAddr)

	demoHost := network.DemoHostCreate(node, host)
	log.Println(demoHost.GetNode().GetIP())

	stdReader := bufio.NewReader(os.Stdin)

	for {
			fmt.Println("\nDemochain (-help to list options):")
			option, err := stdReader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			option = strings.Replace(option, "\n", "", -1)
			option = strings.Replace(option, "\r", "", -1)

			switch option {
				case "-help":
					fmt.Println(" -connect: Connects to Edge Node (only the ELTarget is seted) and open a channel for receive new connections (only the HLNodes is seted)")
					fmt.Println(" -disconnect: Disconnects to Edge Node (only the ELTarget is seted) and close channel for receive new connections (only the HLNodes is seted)")
					fmt.Println(" -add: Adds a new block")
					fmt.Println(" -load: Load blockchain from file")
					fmt.Println(" -clean: Clean blockchain loaded")
					fmt.Println(" -request: Requests blockchain for the Edge Node (only the ELTarget is seted)")

				case "-add":
					fmt.Println("Waiting for data:")
					sendData, err := stdReader.ReadString('\n')
					if err != nil {
						log.Fatal(err)
					}

					sendData = strings.Replace(sendData, "\n", "", -1)
					sendData = strings.Replace(sendData, "\r", "", -1)
					demoHost.ProcessBlock(sendData)

				case "-load":
					demoHost.BlockchainLoad()

				case "-clean":
					demoHost.BlockchainClean()

			  case "-request":
					demoHost.BlockchainRequest()

				case "-connect":
					demoHost.Connect()

				case "-disconnect":
					demoHost.Disconnect()

				default:
					fmt.Println("Invalid option")
			}
	}

	//select {}
}
