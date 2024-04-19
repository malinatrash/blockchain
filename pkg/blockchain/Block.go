package blockchain

type Block struct {
	ID           int64         `json:"id"`
	PreviousHash string        `json:"previousHash"`
	Proof        int64         `json:"proof"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
}
