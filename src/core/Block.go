package core

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
	"fmt"
	"strings"
	"log"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	Data       int
	Hash      string
	PrevHash  string
	Consensus int
	Difficulty int
	Nonce			string
}

// make sure block is valid by checking index, and comparing the hash of the previous block
func IsBlockValid(newBlock, oldBlock Block, consensus int) bool {
	if oldBlock.Index+1 != newBlock.Index {
		log.Println("oldBlock.Index = %s", string(oldBlock.Index))
		log.Println("newBlock.Index = %s", string(newBlock.Index))
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		log.Println("oldBlock.Hash = %s", string(oldBlock.Hash))
		log.Println("newBlock.PrevHash = %s", string(newBlock.PrevHash))
		return false
	}

	if calculateHash(newBlock.Data, newBlock.PrevHash, newBlock.Nonce, consensus) != newBlock.Hash {
		return false
	}

	return true
}

// SHA256 hashing
func calculateHash(data int, prevHash string, nonce string, consensus int) string {

	record := makeRecord(data, prevHash, nonce)

	//POW
		if consensus == 1 {
			h := sha256.New()
			h.Write([]byte(record))
			hashed := h.Sum(nil)
			return hex.EncodeToString(hashed)
		}
		return ""
}

func makeRecord(data int, prevHash string, nonce string) string {
		//return strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.Data) + block.PrevHash + block.Nonce
		return strconv.Itoa(data) + prevHash + nonce
		//return block.Nonce
}

func isHashValid(hash string, difficulty int) bool {
        prefix := strings.Repeat("0", difficulty)
        return strings.HasPrefix(hash, prefix)
}

// create a new block using previous block's hash
func GenerateBlock(oldBlock Block, data int, newHash string, nonce string, consensus int, difficulty int) Block {

	var newBlock Block

	t := time.Now()

	newBlock.Index 		  = oldBlock.Index + 1
	newBlock.Timestamp  = t.String()
	newBlock.Data 		  = data
	newBlock.PrevHash   = oldBlock.Hash
	newBlock.Consensus  = consensus
	newBlock.Difficulty = difficulty

	newBlock.Hash				= newHash
	newBlock.Nonce			= nonce

  return newBlock
}

/*func MinerBlock(data int, consensus int, difficulty int) (string, string) {
	//PoW
	var hash string
	var nonce string
	if consensus == 1 {
		for i := 0; ; i++ {
	          hex := fmt.Sprintf("%x", i)
	          nonce = hex
	          if !isHashValid(calculateHash(data, nonce, consensus), difficulty) {
	                  fmt.Println(calculateHash(data, nonce, consensus), " do more work!")
	                  //time.Sleep(time.Second)
	                  continue
	          } else {
	                  fmt.Println(calculateHash(data, nonce, consensus), " work done!")
	                  hash = calculateHash(data, nonce, consensus)
	                  break
	          }

	  }
	}

	return hash, nonce;
}*/

func MinerBlockLoop(i int, data int, prevHash string, consensus int, difficulty int) (bool, string, string) {
	hex := fmt.Sprintf("%x", i)
	nonce := hex
	if !isHashValid(calculateHash(data, prevHash, nonce, consensus), difficulty) {
					fmt.Println(calculateHash(data, prevHash, nonce, consensus), " do more work!")
					//time.Sleep(time.Second)
					return false, "", ""
	} else {
					fmt.Println(calculateHash(data, prevHash, nonce, consensus), " work done!")
					hash := calculateHash(data, prevHash, nonce, consensus)
					return true, hash, nonce
	}
}


func GenerateGenesisBlock(consensus int, difficulty int) Block {
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateHash(0, "", "0", consensus), "", consensus, difficulty, ""}
	return genesisBlock
}
