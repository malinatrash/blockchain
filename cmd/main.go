package main

import (
	"blockchain/internal/handlers"
	"blockchain/pkg/blockchain"
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
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
	router.GET("/balance", bc.GetBalance)

	router.GET("/wallet", handlers.CreateWallet)
	router.GET("/wallet/all", handlers.GetAllWallets)
	router.GET("/wallet/download", handlers.DownloadPrivateKey)
	router.POST("/wallet/get", handlers.GetWallet)

	router.POST("/nodes/register", func(c *gin.Context) {
		handlers.RegisterNodes(c, bc)
	})

	router.GET("/nodes/resolve", func(c *gin.Context) {
		handlers.ResolveNodes(c, bc)
	})

	for port := 8080; port <= 8090; port++ {
		ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err == nil {
			err := ln.Close()
			if err != nil {
				return
			}
			err = router.Run(":" + strconv.Itoa(port))
			if err != nil {
				continue
			}
			break
		}
	}
}
