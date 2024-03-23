package blockchain

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type Wallet struct {
	PrivateKey *rsa.PrivateKey `json:"privateKey"`
	PublicKey  *rsa.PublicKey  `json:"publicKey"`
	Address    string          `json:"address"`
	TimeStamp  string          `json:"timeStamp"`
}

var Wallets map[string]*Wallet

func init() {
	Wallets = make(map[string]*Wallet)
}

func GenerateWallet() (*Wallet, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	publicKey := &privateKey.PublicKey
	address := GenerateAddress(publicKey)
	wallet := &Wallet{PrivateKey: privateKey, PublicKey: publicKey, Address: address, TimeStamp: time.Now().Format(time.RFC3339)}
	Wallets[address] = wallet
	fmt.Printf("Wallet %s created\n", address)
	return wallet, nil
}

func GenerateAddress(publicKey *rsa.PublicKey) string {
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	hash := sha256.Sum256(publicKeyBytes)
	return hex.EncodeToString(hash[:])
}

func GetWallet(address string) (*Wallet, error) {
	wallet, ok := Wallets[address]
	if !ok {
		return nil, errors.New("wallet not found")
	}
	return wallet, nil
}
