package blockchain

type Block struct {
	id           int64
	timestamp    string
	transactions []Transaction
	proof        int64
	previousHash string
}
