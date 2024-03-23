package handlers

import (
	"blockchain/pkg/blockchain"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)

func CreateWallet(c *gin.Context) {
	wallet, err := blockchain.GenerateWallet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address": wallet.Address,
	})
}

func DownloadPrivateKey(c *gin.Context) {
	address := c.Query("address")
	wallet, err := blockchain.GetWallet(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	file, err := createPrivateKeyFile(wallet.PrivateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer func() {
		if err := os.Remove(file.Name()); err != nil {
		}
	}()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=private_key.pem")
	c.Header("Content-Type", "application/octet-stream")

	c.File(file.Name())
}

func createPrivateKeyFile(privateKey *rsa.PrivateKey) (*os.File, error) {
	privateKeyFileName := "private_key.pem"
	filePath := filepath.Join(".", privateKeyFileName)
	if privateKey == nil {
		return nil, errors.New("private key is nil")
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
		}
	}()

	err = pem.Encode(file, block)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func GetWallet(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()

	privateKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PEM file"})
		return
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse private key"})
		return
	}

	address := GetWalletAddress(privateKey)
	if address == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"address": address})
}

func GetWalletAddress(privateKey *rsa.PrivateKey) string {
	for _, wallet := range blockchain.Wallets {
		if reflect.DeepEqual(privateKey, wallet.PrivateKey) {
			return wallet.Address
		}
	}
	return ""
}

func GetAllWallets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"wallets": blockchain.Wallets,
	})
}
