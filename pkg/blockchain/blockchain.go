package blockchain

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Blockchain struct {
	Chain              []Block
	PoolOfTransactions []Transaction
	mu                 sync.Mutex
}

func NewBlockchain() *Blockchain {
	blockchain := &Blockchain{
		Chain:              []Block{},
		PoolOfTransactions: []Transaction{},
	}
	blockchain.NewBlock(100)
	return blockchain
}

func (bc *Blockchain) NewBlock(proof int64) Block {
	return bc.newBlock(proof, "")
}

func (bc *Blockchain) LastBlock() *Block {
	return &bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) NewBlockWithPreviousHash(proof int64, previousHash string) Block {
	return bc.newBlock(proof, previousHash)
}

func (bc *Blockchain) newBlock(proof int64, previousHash string) Block {
	block := Block{
		id:           int64(len(bc.Chain) + 1),
		timestamp:    time.Now().Format(time.RFC3339),
		transactions: bc.PoolOfTransactions,
		proof:        proof,
		previousHash: previousHash,
	}
	bc.Chain = append(bc.Chain, block)
	bc.PoolOfTransactions = []Transaction{}
	return block
}

func (bc *Blockchain) NewTransaction(amount int64, recipient string, sender string) (*int64, error) {
	if sender == "0" {
		bc.PoolOfTransactions = append(bc.PoolOfTransactions, Transaction{
			Sender:    sender,
			Recipient: recipient,
			Amount:    amount,
			signature: "",
		})
	} else {
		senderWallet, err := GetWallet(sender)
		if err != nil {
			return nil, fmt.Errorf("Sender does not exists")
		}
		_, err = GetWallet(recipient)
		if err != nil {
			return nil, fmt.Errorf("Recipient does not exists")
		}
		balance, err := bc.Balance(sender)
		if err != nil {
			return nil, fmt.Errorf("Sender does not exists")
		}
		if *balance < amount {
			return nil, fmt.Errorf("Sender balance is lower than amount")
		}

		message := fmt.Sprintf("%d%s%s", amount, recipient, senderWallet.Address)
		hashed := sha256.Sum256([]byte(message))

		signature, err := rsa.SignPKCS1v15(rand.Reader, senderWallet.PrivateKey, crypto.SHA256, hashed[:])
		if err != nil {
			return nil, fmt.Errorf("Failed to sign transaction")
		}

		bc.PoolOfTransactions = append(bc.PoolOfTransactions, Transaction{
			Sender:    senderWallet.Address,
			Recipient: recipient,
			Amount:    amount,
			signature: string(signature),
		})
	}

	lastTransactionId := bc.LastBlock()
	newTransactionId := lastTransactionId.id + 1
	return &newTransactionId, nil
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
				transactionString += tx.Sender + tx.Recipient + strconv.FormatInt(tx.Amount, 10)
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
	return guessHashStr[:6] == "000000"
}
