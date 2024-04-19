package handlers

import (
	"blockchain/pkg/blockchain"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterNodes(context *gin.Context, bc *blockchain.Blockchain) {
	var request struct {
		Nodes []string `json:"nodes"`
	}
	if err := context.BindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, node := range request.Nodes {
		err := bc.RegisterNode(node)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	response := gin.H{
		"message":     "New nodes have been added",
		"total_nodes": request.Nodes,
	}
	context.JSON(http.StatusCreated, response)
}

func ResolveNodes(context *gin.Context, bc *blockchain.Blockchain) {
	if bc.ResolveConflicts() {
		context.JSON(http.StatusOK, gin.H{"message": "Our chain was replaced", "new_chain": bc.Chain})
	} else {
		context.JSON(http.StatusOK, gin.H{"message": "Our chain is authoritative", "chain": bc.Chain})
	}
}
