package blockchain

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (bc *Blockchain) GetNewTransaction(c *gin.Context) {
	var transaction Transaction
	if err := c.BindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	index := bc.NewTransaction(transaction.amount, transaction.recipient, transaction.sender)
	message := fmt.Sprintf("Transaction added to the block %d", index)
	c.JSON(http.StatusCreated, gin.H{"message": message})
}

func (bc *Blockchain) Mine(c *gin.Context) {
	lastBlock := bc.LastBlock()
	lastProof := lastBlock.proof
	proof := bc.ProofOfWork(lastProof)
	bc.NewTransaction(1, uuid.New().String(), "0")

	previousHash := bc.Hash(*lastBlock)
	block := bc.newBlock(proof, previousHash)

	var transactions []gin.H
	for _, transaction := range block.transactions {
		transactions = append(transactions, gin.H{
			"amount":    transaction.amount,
			"recipient": transaction.recipient,
			"sender":    transaction.sender,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "New Block Forged",
		"index":        block.id,
		"transactions": transactions,
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
				"amount":    transaction.amount,
				"recipient": transaction.recipient,
				"sender":    transaction.sender,
			})
		}
		blocks = append(blocks, gin.H{
			"id":           block.id,
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
