package core

import (
	protocol "github.com/libp2p/go-libp2p-protocol"
	crypto "github.com/libp2p/go-libp2p-crypto"
	"log"
	"io/ioutil"
	"strconv"
	"strings"
)

type hlNode struct {
	HLTarget string
	Permiss int //0 - SOMENTE LEITURA 1 - LEITURA E ESCRITA
}

type Node struct {
	//Consensus	Consensus //CONTEM CONFIGUCOES QUE FORAM UTILIZADAS
	IP string
	Port int
	NetworkName protocol.ID
	PublicKey crypto.PubKey
	PrivateKey crypto.PrivKey
	CryptographicType int
  CryptographicBits int
	ELTarget string
	PathBlockchainFile string
	Permissioned int //0 NAO e 1 SIM
	HLNodes []hlNode
	Consensus int //1 pow | 2 pos | 3 pbft | 4 bftraft | 5 ripple
}

func nodeObjectCreate(ip string, port int, networkName protocol.ID, publicKey crypto.PubKey, privateKey crypto.PrivKey, cryptographicType int, cryptographicBits int, elTarget string, pathBlockchainFile string, hlNodes []hlNode, consensus int) (Node) {
	var node Node
	node.IP 							  = ip
	node.Port 						  = port
	node.NetworkName 			  = networkName
	node.PublicKey 				  = publicKey
	node.PrivateKey 			  = privateKey
	node.CryptographicType  = cryptographicType
  node.CryptographicBits  = cryptographicBits
	node.ELTarget 				  = elTarget
	node.PathBlockchainFile = pathBlockchainFile
	node.HLNodes 						= hlNodes
	node.Consensus					= consensus

	return node
}

func NodeLoad(ip string, port string, networkName string, pathKey string, cryptographicType string, cryptographicBits string, elTarget string, pathBlockchainFile string, hlNodes string, consensus string) (Node) {
	portInt, _ 							:= strconv.Atoi(port)
	cryptographicTypeInt, _ := strconv.Atoi(cryptographicType)
	cryptographicBitsInt, _ := strconv.Atoi(cryptographicBits)
	consensusInt, _ := strconv.Atoi(consensus)

	pub, priv := getPrivateKey(pathKey, cryptographicTypeInt, cryptographicBitsInt)

	return nodeObjectCreate(ip, portInt, protocol.ID(networkName), pub, priv, cryptographicTypeInt, cryptographicBitsInt, elTarget, pathBlockchainFile, formatHLNodes(hlNodes), consensusInt)
}

func formatHLNodes(hlNodes string) ([] hlNode) {
	var hlNodesObject []hlNode

	if hlNodes != "" {
		hlNodesArray := strings.Split(hlNodes, "|")
		log.Println(hlNodesArray)

		var hlNodeI hlNode

		for _, hlNodeIndex := range hlNodesArray {
				hlNodeD := strings.Split(hlNodeIndex, "-") //PEGA O RESUMO E A PERMISSAO
				hlNodeI.HLTarget = hlNodeD[0]
				permiss, _ := strconv.Atoi(hlNodeD[1])
				hlNodeI.Permiss	= permiss
				hlNodesObject = append(hlNodesObject, hlNodeI)
		}
	}

	return hlNodesObject
}


func getPrivateKey(path string, cryptographicType int, cryptographicBits int) (crypto.PubKey, crypto.PrivKey) {
	log.Println("Buscando private key")
	var priv crypto.PrivKey
	pub, priv, fileExists := readFilePrivateKey(path)

	if fileExists == false {
		log.Println("Gerando nova private key")
		pub_nova, priv_nova := writeFilePrivateKey(path, cryptographicType, cryptographicBits)
		priv = priv_nova
		pub = pub_nova
	}

	return pub, priv;
}

func writeFilePrivateKey(path string, cryptographicType int, cryptographicBits int) (crypto.PubKey, crypto.PrivKey) {
	priv, pub, err := crypto.GenerateKeyPair(cryptographicType, cryptographicBits)

	bytes_pk, _ := crypto.MarshalPrivateKey(priv)
	err = ioutil.WriteFile(path, bytes_pk, 0644)
  if err != nil {
  	log.Fatal(err)
  }

	return pub, priv;
}

func readFilePrivateKey(path string) (crypto.PubKey, crypto.PrivKey, bool)  {
	filePrivKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, false
	}

	privKey, _ := crypto.UnmarshalPrivateKey(filePrivKey)
	return privKey.GetPublic(), privKey, true
}
