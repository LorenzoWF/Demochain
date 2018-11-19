package main

import (
	"flag"
	"log"
	"os"
  "io"
	golog "github.com/ipfs/go-log"
	gologging "github.com/whyrusleeping/go-logging"
	"github.com/joho/godotenv"
	"bufio"
	"strings"
	//"fmt"
  "encoding/csv"
	core "demochain/core"
	network "demochain/network"
  "database/sql"
)

const (
  host_db = "localhost"
  port_db = 5432
  user_db = "postgres"
  pass_db = "Post@sol"
  name_db = "scenario2"
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

	err = godotenv.Load(path)
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

	if err != nil {
		log.Fatal(err)
	}

	log.Println(node.GetTarget())


	fullAddr := network.MakeFullAddr(host)
	log.Println(fullAddr)

	demoHost := network.DemoHostCreate(node, host)
	log.Println(demoHost.GetNode().GetIP())

	demoHost.BlockchainLoad()
	demoHost.Connect()

	stdReader := bufio.NewReader(os.Stdin)

  option, err := stdReader.ReadString('\n')
  if err != nil {
    log.Fatal(err)
  }
  option = strings.Replace(option, "\n", "", -1)
  option = strings.Replace(option, "\r", "", -1)

  if option == "-execute" {
    csvFile, _ := os.Open("beach-water-quality-automated-sensors-1.csv")
    reader := csv.NewReader(bufio.NewReader(csvFile))

    firstLine := true

    for {
       line, error := reader.Read()

       if firstLine == true {
         firstLine = false
         continue
       }

       if error == io.EOF {
           break
       } else if error != nil {
           log.Fatal(error)
       }

      lineFull := line[0] + "|" + line[1] + "|" + line[2] + "|" + line[3] + "|" + line[4] + "|" + line[5] + "|" + line[6] + "|" + line[7] + "|" + line[8]
      log.Println(lineFull)
      demoHost.ProcessBlock(lineFull)
    }
  }

	//select {}
}
