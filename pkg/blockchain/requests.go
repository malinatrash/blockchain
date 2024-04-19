package blockchain

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
	lastProof := lastBlock.Proof
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

	block.Transactions = append(block.Transactions, transactions...)

	c.JSON(http.StatusOK, gin.H{
		"message":      "New Block Forged",
		"Id":           block.ID,
		"Timestamp":    block.Timestamp,
		"Transactions": block.Transactions,
		"Proof":        block.Proof,
		"PreviousHash": block.PreviousHash,
	})
}

func (bc *Blockchain) GetChain(c *gin.Context) {
	var blocks []gin.H
	for _, block := range bc.Chain {
		var transactions []gin.H
		for _, transaction := range block.Transactions {
			transactions = append(transactions, gin.H{
				"amount":    transaction.Amount,
				"recipient": transaction.Recipient,
				"sender":    transaction.Sender,
			})
		}
		blocks = append(blocks, gin.H{
			"Id":           block.ID,
			"Timestamp":    block.Timestamp,
			"Transactions": transactions,
			"Proof":        block.Proof,
			"PreviousHash": block.PreviousHash,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"chain":  blocks,
		"length": len(bc.Chain),
	})
}

func (bc *Blockchain) Balance(address string) (*int64, error) {
	if Wallets[address] == nil {
		return nil, fmt.Errorf("Wallet is not exists")
	}
	balance := int64(0)

	for _, block := range bc.Chain {
		for _, transaction := range block.Transactions {
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
