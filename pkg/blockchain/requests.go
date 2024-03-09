package blockchain

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (bc *Blockchain) GetNewTransaction(c *gin.Context) {
	var transaction Transaction
	if err := c.BindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bc.CurrentTransactions = append(bc.CurrentTransactions, transaction)
	c.JSON(http.StatusCreated, gin.H{"message": "Transaction added to the block"})
}

func (bc *Blockchain) Mine(c *gin.Context) {
	lastBlock := bc.LastBlock()
	lastProof := lastBlock.proof
	proof := bc.ProofOfWork(lastProof)

	bc.GetNewTransaction(c)

	previousHash := bc.Hash(*lastBlock)
	block := Block{
		id:           lastBlock.id + 1,
		timestamp:    time.Now().Format(time.RFC3339),
		transactions: bc.CurrentTransactions,
		proof:        proof,
		previousHash: previousHash,
	}

	bc.CurrentTransactions = nil
	bc.Chain = append(bc.Chain, block)

	c.JSON(http.StatusOK, gin.H{
		"message":      "New Block Forged",
		"index":        block.id,
		"transactions": block.transactions,
		"proof":        block.proof,
		"previousHash": block.previousHash,
	})
}

func (bc *Blockchain) GetChain(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"chain":  bc.Chain,
		"length": len(bc.Chain),
	})
}
