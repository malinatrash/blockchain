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
	fmt.Printf("%v", transaction)
	fmt.Printf("%v", "!!!!!!!")
	index := bc.NewTransaction(transaction.Amount, transaction.Recipient, transaction.Sender)
	message := fmt.Sprintf("Transaction added to the block %d", index)
	c.JSON(http.StatusCreated, gin.H{"message": message})
}

func (bc *Blockchain) Mine(c *gin.Context) {
	lastBlock := bc.LastBlock()
	lastProof := lastBlock.proof
	proof := bc.ProofOfWork(lastProof)
	bc.NewTransaction(1, NodeIdentifier, "0")

	previousHash := bc.Hash(*lastBlock)
	block := bc.newBlock(proof, previousHash)

	var transactions []gin.H
	for _, transaction := range block.transactions {
		transactions = append(transactions, gin.H{
			"amount":    transaction.Amount,
			"recipient": transaction.Recipient,
			"sender":    transaction.Sender,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "New Block Forged",
		"id":           block.id,
		"timestamp":    block.timestamp,
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
