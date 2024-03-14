package main

import (
	"blockchain/pkg/blockchain"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	bc := blockchain.NewBlockchain()

	router.POST("/transactions/new", bc.GetNewTransaction)
	router.GET("/mine", bc.Mine)
	router.GET("/chain", bc.GetChain)

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
