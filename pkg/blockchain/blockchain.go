package blockchain

import (
	"blockchain/pkg"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type Blockchain struct {
	Chain              []Block
	PoolOfTransactions []Transaction
	Nodes              pkg.Set
}

func NewBlockchain() *Blockchain {
	blockchain := &Blockchain{
		Chain:              []Block{},
		PoolOfTransactions: []Transaction{},
		Nodes:              make(pkg.Set),
	}
	blockchain.NewBlock(100)
	return blockchain
}

func (bc *Blockchain) RegisterNode(address string) error {
	parsedURL, err := url.Parse(address)
	if err != nil {
		return err
	}
	if parsedURL.Host != "" {
		bc.Nodes.Add(parsedURL.Host)
	} else if parsedURL.Path != "" {
		bc.Nodes.Add(parsedURL.Path)
	} else {
		return errors.New("invalid URL")
	}
	return nil
}

func (bc *Blockchain) ValidChain(chain []Block) bool {
	lastBlock := chain[0]
	currentIndex := 1

	for currentIndex < len(chain) {
		block := chain[currentIndex]

		lastBlockHash := bc.Hash(lastBlock)
		if block.PreviousHash != lastBlockHash {
			return false
		}

		if !bc.ValidProof(lastBlock.Proof, block.Proof) {
			return false
		}

		lastBlock = block
		currentIndex++
	}

	return true
}

func (bc *Blockchain) ResolveConflicts() bool {
	newChain := make([]Block, 0)

	maxLength := len(bc.Chain)
	for node := range bc.Nodes {
		response, err := http.Get(fmt.Sprintf("http://%s/chain", node))
		if err != nil {
			continue
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			var data struct {
				Chain  []Block `json:"chain"`
				Length int     `json:"length"`
			}
			if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
				fmt.Println("Error decoding JSON:", err)
				return false
			}

			if data.Length > maxLength {
				maxLength = data.Length
				newChain = data.Chain
			}
		}
	}

	if len(newChain) > 0 {
		bc.Chain = newChain
		return true
	}

	return false
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
		ID:           int64(len(bc.Chain) + 1),
		Timestamp:    time.Now().Format(time.RFC3339),
		Transactions: bc.PoolOfTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
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
			return nil, fmt.Errorf("sender does not exists")
		}
		_, err = GetWallet(recipient)
		if err != nil {
			return nil, fmt.Errorf("recipient does not exists")
		}
		balance, err := bc.Balance(sender)
		if err != nil {
			return nil, fmt.Errorf("sender does not exists")
		}
		if *balance < amount {
			return nil, fmt.Errorf("sender balance is lower than amount")
		}

		message := fmt.Sprintf("%d%s%s", amount, recipient, senderWallet.Address)
		hashed := sha256.Sum256([]byte(message))

		signature, err := rsa.SignPKCS1v15(rand.Reader, senderWallet.PrivateKey, crypto.SHA256, hashed[:])
		if err != nil {
			return nil, fmt.Errorf("failed to sign transaction")
		}

		bc.PoolOfTransactions = append(bc.PoolOfTransactions, Transaction{
			Sender:    senderWallet.Address,
			Recipient: recipient,
			Amount:    amount,
			signature: string(signature),
		})
	}

	lastTransactionId := bc.LastBlock()
	newTransactionId := lastTransactionId.ID + 1
	return &newTransactionId, nil
}

func (bc *Blockchain) Hash(block Block) string {
	blockJSON, _ := json.Marshal(block)
	var keys []string
	keys = append(keys, "Id", "Timestamp", "Transactions", "Proof", "PreviousHash")
	sort.Strings(keys)

	blockString := ""
	for _, k := range keys {
		switch k {
		case "Transactions":
			transactionString := ""
			for _, tx := range block.Transactions {
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
