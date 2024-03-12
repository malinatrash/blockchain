package main

import (
	"blockchain/pkg/blockchain"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	bc := blockchain.NewBlockchain()

	router.POST("/transactions/new", bc.GetNewTransaction)
	router.GET("/mine", bc.Mine)
	router.GET("/chain", bc.GetChain)

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
