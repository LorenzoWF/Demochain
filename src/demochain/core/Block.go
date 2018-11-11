package core

import (
	"crypto/sha256"
	"encoding/hex"
	//"strconv"
	"time"
	"fmt"
	"strings"
	//"log"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      string
	PrevHash  string
	Consensus Consensus
	Nonce			string
	Target	  string
}

func (block Block) getBlockPublic() {

}

func (block Block) GetIndex() (int) {
		return block.Index
}

func (block Block) GetTimestamp() (string) {
		return block.Timestamp
}

func (block Block) GetData() (string) {
		return block.Data
}

func (block Block) GetHash() (string) {
		return block.Hash
}

func (block Block) GetPrevHash() (string) {
		return block.PrevHash
}

func (block Block) GetConsensus() (Consensus) {
		return block.Consensus
}

func (block Block) GetNonce() (string) {
		return block.Nonce
}

func (block Block) GetTarget() (string) {
		return block.Target
}

// make sure block is valid by checking index, and comparing the hash of the previous block
func (newBlock Block) IsBlockValid(oldBlock Block) bool {
	if oldBlock.GetIndex()+1 != newBlock.GetIndex() {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "Rejected block! Invalid Index:")
		fmt.Printf("\x1b[31m%s%d\x1b[0m\n", "Old block Index + 1: ", oldBlock.GetIndex()+1)
		fmt.Printf("\x1b[31m%s%d\x1b[0m\n", "New block Index: ", newBlock.GetIndex())
		//log.Println("oldBlock.Index+1 = ", oldBlock.Index+1)
		//log.Println("newBlock.Index = ", newBlock.Index)
		return false
	}

	if oldBlock.GetConsensus() != newBlock.GetConsensus() {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "Rejected block! Invalid Consensus:")
		fmt.Printf("\x1b[31m%s%d\x1b[0m\n", "Old block consensus type: ", oldBlock.GetConsensus().GetType())
		fmt.Printf("\x1b[31m%s%d\x1b[0m\n", "Old block consensus difficulty: ", oldBlock.GetConsensus().GetDifficulty())
		fmt.Printf("\x1b[31m%s%d\x1b[0m\n", "New block consensus type: ", newBlock.GetConsensus().GetType())
		fmt.Printf("\x1b[31m%s%d\x1b[0m\n", "New block consensus difficulty: ", newBlock.GetConsensus().GetDifficulty())
		//log.Println("oldBlock.Consensus = ", oldBlock.Consensus)
		//log.Println("newBlock.Consensus = ", newBlock.Consensus)
		return false
	}

	if oldBlock.GetHash() != newBlock.GetPrevHash() {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "Rejected block! Invalid Previous Hash:")
		fmt.Printf("\x1b[31m%s%s\x1b[0m\n", "Old block hash: ", oldBlock.GetHash())
		fmt.Printf("\x1b[31m%s%s\x1b[0m\n", "New block prev hash: ", newBlock.GetPrevHash())
		//log.Println("oldBlock.Hash = ", oldBlock.Hash)
		//log.Println("newBlock.Hash = ", newBlock.Hash)
		return false
	}

	calculatedHash := calculateHash(newBlock.GetData(), newBlock.GetPrevHash(), newBlock.GetNonce(), newBlock.GetConsensus().GetType())
	if calculatedHash != newBlock.GetHash() {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "Rejected block! Invalid Hash:")
		fmt.Printf("\x1b[31m%s%s\x1b[0m\n", "Calculated hash: ", calculatedHash)
		fmt.Printf("\x1b[31m%s%s\x1b[0m\n", "New block hash: ", newBlock.GetHash())
		//log.Println("func = ", calculateHash(newBlock.Data, newBlock.PrevHash, newBlock.Nonce, newBlock.Consensus.Type))
		//log.Println("newBlock.Hash = ", newBlock.Hash)
		return false
	}

	return true
}

// SHA256 hashing
func calculateHash(data string, prevHash string, nonce string, consensusType int) string {

	record := makeRecord(data, prevHash, nonce)

	//POW
		if consensusType == 1 {
			h := sha256.New()
			h.Write([]byte(record))
			hashed := h.Sum(nil)
			return hex.EncodeToString(hashed)
		}
		return ""
}

func makeRecord(data string, prevHash string, nonce string) string {
		return data + prevHash + nonce
}

func isHashValid(hash string, difficulty int) bool {
        prefix := strings.Repeat("0", difficulty)
        return strings.HasPrefix(hash, prefix)
}

// create a new block using previous block's hash
func GenerateBlock(oldBlock Block, data string, newHash string, nonce string, consensus Consensus, target string) Block {
	t := time.Now()
	newBlock := Block{}
	newBlock = Block{oldBlock.GetIndex() + 1, t.String(), data, newHash, oldBlock.GetHash(), consensus, nonce, target}

  return newBlock
}

//Proof of Work
func MinerBlockLoop(i int, data string, prevHash string, consensus Consensus) (bool, string, string) {
	hex := fmt.Sprintf("%x", i)
	nonce := hex
	if !isHashValid(calculateHash(data, prevHash, nonce, consensus.GetType()), consensus.GetDifficulty()) {
					//fmt.Println(calculateHash(data, prevHash, nonce, consensus.Type), " do more work!")
					//time.Sleep(time.Second)
					return false, "", ""
	} else {
					//fmt.Println(calculateHash(data, prevHash, nonce, consensus.Type), " work done!")
					hash := calculateHash(data, prevHash, nonce, consensus.GetType())
					return true, hash, nonce
	}
}


func GenerateGenesisBlock(consensus Consensus, target string) Block {
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), "0", calculateHash("0", "", "0", consensus.GetType()), "", consensus, "", target}
	return genesisBlock
}
