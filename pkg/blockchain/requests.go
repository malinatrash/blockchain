package blockchain

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func (bc *Blockchain) GetNewTransaction(c *gin.Context) {
	var transaction Transaction
	if err := c.BindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	index, err := bc.NewTransaction(transaction.Amount, transaction.Recipient, transaction.Sender)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message := fmt.Sprintf("Transaction added to the block %d", index)
	c.JSON(http.StatusCreated, gin.H{"message": message})
}

func (bc *Blockchain) Mine(c *gin.Context) {
	var Miner = c.Query("address")
	if Miner == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid address",
		})
		return
	}
	_, err := bc.NewTransaction(1, Miner, "0")
	if err != nil {
		return
	}

	lastBlock := bc.LastBlock()
	lastProof := lastBlock.proof
	proof := bc.ProofOfWork(lastProof)

	previousHash := bc.Hash(*lastBlock)

	block := bc.newBlock(proof, previousHash)

	var transactions []Transaction
	for _, transaction := range bc.PoolOfTransactions {
		tx := Transaction{
			Amount:    transaction.Amount,
			Recipient: transaction.Recipient,
			Sender:    transaction.Sender,
		}
		transactions = append(transactions, tx)
	}

	block.transactions = append(block.transactions, transactions...)

	c.JSON(http.StatusOK, gin.H{
		"message":      "New Block Forged",
		"id":           block.id,
		"timestamp":    block.timestamp,
		"transactions": block.transactions,
		"proof":        block.proof,
		"previousHash": block.previousHash,
	})
}

func (bc *Blockchain) GetChain(c *gin.Context) {
	var blocks []gin.H
	for _, block := range bc.Chain {
		var transactions []gin.H
		for _, transaction := range block.transactions {
			transactions = append(transactions, gin.H{
				"amount":    transaction.Amount,
				"recipient": transaction.Recipient,
				"sender":    transaction.Sender,
			})
		}
		blocks = append(blocks, gin.H{
			"id":           block.id,
			"timestamp":    block.timestamp,
			"transactions": transactions,
			"proof":        block.proof,
			"previousHash": block.previousHash,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"chain":  blocks,
		"length": len(bc.Chain),
	})
}

func (bc *Blockchain) CreateWallet(c *gin.Context) {
	wallet, err := GenerateWallet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	//privateKeyFileName := "private_key.pem"
	//err = WritePrivateKeyToFile(wallet.PrivateKey, privateKeyFileName)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{})
	//	return
	//}

	//c.Header("Content-Description", "File Transfer")
	//c.Header("Content-Disposition", "attachment; filename="+privateKeyFileName)
	//c.Header("Content-Type", "application/octet-stream")
	//c.Header("Content-Transfer-Encoding", "binary")
	//c.Header("Expires", "0")
	//c.Header("Cache-Control", "must-revalidate")
	//c.Header("Pragma", "public")

	c.JSON(http.StatusOK, gin.H{
		"address":    wallet.Address,
		"privateKey": x509.MarshalPKCS1PrivateKey(wallet.PrivateKey),
	})
	//filePath := filepath.Join(".", privateKeyFileName)
	//c.File(filePath)
}

func WritePrivateKeyToFile(privateKey *rsa.PrivateKey, fileName string) error {
	if privateKey == nil {
		return errors.New("private key is nil")
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	filePath := filepath.Join(".", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	return nil
}

func (bc *Blockchain) Balance(address string) (*int64, error) {
	if Wallets[address] == nil {
		return nil, fmt.Errorf("Wallet is not exists")
	}
	balance := int64(0)

	for _, block := range bc.Chain {
		for _, transaction := range block.transactions {
			if transaction.Sender == address {
				balance -= transaction.Amount
			}
			if transaction.Recipient == address {
				balance += transaction.Amount
			}
		}
	}
	return &balance, nil
}

func (bc *Blockchain) GetBalance(c *gin.Context) {
	var address = c.Query("address")
	balance, err := bc.Balance(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wallet with this address does not exists",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"balance": balance,
		})
	}
}
