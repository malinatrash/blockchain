package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
)

type Blockchain struct {
	Chain               []Block
	CurrentTransactions []Transaction
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain:               make([]Block, 0),
		CurrentTransactions: make([]Transaction, 0),
	}
}

func (bc *Blockchain) NewBlock(proof int64) Block {
	return bc.newBlock(proof, "")
}

func (bc *Blockchain) LastBlock() *Block {
	if len(bc.Chain) == 0 {
		return nil
	}
	return &bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) NewBlockWithPreviousHash(proof int64, previousHash string) Block {
	return bc.newBlock(proof, previousHash)
}

func (bc *Blockchain) newBlock(proof int64, previousHash string) Block {
	block := Block{
		id:           int64(len(bc.Chain) + 1),
		timestamp:    time.Now().Format(time.RFC3339),
		transactions: bc.CurrentTransactions,
		proof:        proof,
		previousHash: previousHash,
	}

	bc.CurrentTransactions = nil

	bc.Chain = append(bc.Chain, block)

	return block
}

func (bc *Blockchain) NewTransaction(sender string, recipient string, amount int64) int64 {
	bc.CurrentTransactions = append(bc.CurrentTransactions, Transaction{
		sender:    sender,
		recipient: recipient,
		amount:    amount,
	})

	return bc.LastBlock().id + 1
}

func (bc *Blockchain) Hash(block Block) string {
	blockJSON, _ := json.Marshal(block)
	var keys []string
	keys = append(keys, "id", "timestamp", "transactions", "proof", "previousHash")
	sort.Strings(keys)

	blockString := ""
	for _, k := range keys {
		switch k {
		case "transactions":
			transactionString := ""
			for _, tx := range block.transactions {
				transactionString += tx.sender + tx.recipient + strconv.FormatInt(tx.amount, 10)
			}
			blockString += transactionString
		default:
			blockString += k + string(blockJSON)
		}
	}
	hash := sha256.Sum256([]byte(blockString))
	return hex.EncodeToString(hash[:])
}

func (bc *Blockchain) ProofOfWork(lastProof int64) int64 {
	proof := int64(0)
	for !bc.ValidProof(lastProof, proof) {
		proof++
	}
	return proof
}

func (bc *Blockchain) ValidProof(lastProof, proof int64) bool {
	guess := fmt.Sprintf("%d %d", lastProof, proof)
	guessHash := sha256.Sum256([]byte(guess))
	guessHashStr := hex.EncodeToString(guessHash[:])
	return guessHashStr[:4] == "0000"
}
