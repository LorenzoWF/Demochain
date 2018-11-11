package core

import (
	protocol "github.com/libp2p/go-libp2p-protocol"
	crypto "github.com/libp2p/go-libp2p-crypto"
	"log"
	"io/ioutil"
	"strconv"
	"strings"
)

type Node struct {
	//Consensus	Consensus //CONTEM CONFIGUCOES QUE FORAM UTILIZADAS
	IP string
	Port int
	NetworkName protocol.ID
	PublicKey crypto.PubKey
	PrivateKey crypto.PrivKey
	Target string
	CryptographicType int
  CryptographicBits int
	ELTarget string
	PathBlockchainFile string
	Permissioned int //0 NAO e 1 SIM
	HLNodes []HLNode
	Consensus Consensus
}

func nodeObjectCreate(ip string, port int, networkName protocol.ID, publicKey crypto.PubKey, privateKey crypto.PrivKey, cryptographicType int, cryptographicBits int, elTarget string, pathBlockchainFile string, hlNodes []HLNode, consensus Consensus) (Node) {
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

func NodeLoad(ip string, port string, networkName string, pathKey string, cryptographicType string, cryptographicBits string, elTarget string, pathBlockchainFile string, hlNodes string, consensusType string, difficulty string) (Node) {
	portInt, _ 							:= strconv.Atoi(port)
	cryptographicTypeInt, _ := strconv.Atoi(cryptographicType)
	cryptographicBitsInt, _ := strconv.Atoi(cryptographicBits)
	consensusTypeInt, _ 		:= strconv.Atoi(consensusType)
	difficultyInt, _ 				:= strconv.Atoi(difficulty)

	consensus := ConsensusCreate(consensusTypeInt, difficultyInt)

	pub, priv := getKeys(pathKey, cryptographicTypeInt, cryptographicBitsInt)

	return nodeObjectCreate(ip, portInt, protocol.ID(networkName), pub, priv, cryptographicTypeInt, cryptographicBitsInt, elTarget, pathBlockchainFile, formatHLNodes(hlNodes), consensus)

}

func formatHLNodes(hlNodes string) ([] HLNode) {
	var hlNodesObject []HLNode

	if hlNodes != "" {
		hlNodesArray := strings.Split(hlNodes, "|")
		log.Println(hlNodesArray)

		for _, hlNodeIndex := range hlNodesArray {
				hlNodeD := strings.Split(hlNodeIndex, "-") //PEGA O RESUMO E A PERMISSAO
				permiss, _ := strconv.Atoi(hlNodeD[1])

				hlNodesObject = append(hlNodesObject, HLNodeCreate(hlNodeD[0], permiss))
		}
	}

	return hlNodesObject
}


func getKeys(path string, cryptographicType int, cryptographicBits int) (crypto.PubKey, crypto.PrivKey) {
	log.Println("Buscando private key")
	var priv crypto.PrivKey
	pub, priv, fileExists := readFileKeys(path)

	if fileExists == false {
		log.Println("Gerando nova private key")
		pub_nova, priv_nova := writeFileKeys(path, cryptographicType, cryptographicBits)
		priv = priv_nova
		pub = pub_nova
	}

	return pub, priv;
}

func writeFileKeys(path string, cryptographicType int, cryptographicBits int) (crypto.PubKey, crypto.PrivKey) {
	priv, pub, err := crypto.GenerateKeyPair(cryptographicType, cryptographicBits)

	bytes_pk, _ := crypto.MarshalPrivateKey(priv)
	err = ioutil.WriteFile(path, bytes_pk, 0644)
  if err != nil {
  	log.Fatal(err)
  }

	return pub, priv;
}

func readFileKeys(path string) (crypto.PubKey, crypto.PrivKey, bool)  {
	filePrivKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, false
	}

	privKey, _ := crypto.UnmarshalPrivateKey(filePrivKey)
	return privKey.GetPublic(), privKey, true
}

func (node Node) GetIP() (string) {
		return node.IP
}

func (node Node) GetPort() (int) {
		return node.Port
}

func (node Node) GetNetworkName() (protocol.ID) {
		return node.NetworkName
}

func (node Node) GetPublicKey() (crypto.PubKey) {
		return node.PublicKey
}

func (node Node) GetPrivateKey() (crypto.PrivKey) {
		return node.PrivateKey
}

func (node Node) GetTarget() (string) {
		return node.Target
}

func (node *Node) SetTarget(target string) {
		node.Target = target
}

func (node Node) GetCryptographicType() (int) {
		return node.CryptographicType
}

func (node Node) GetCryptographicBits() (int) {
		return node.CryptographicBits
}

func (node Node) GetELTarget() (string) {
		return node.ELTarget
}

func (node Node) GetPathBlockchainFile() (string) {
		return node.PathBlockchainFile
}

func (node Node) GetPermissioned() (int) {
		return node.Permissioned
}

func (node Node) GetHLNodes() ([]HLNode) {
		return node.HLNodes
}

func (node Node) GetConsensus() (Consensus) {
		return node.Consensus
}
