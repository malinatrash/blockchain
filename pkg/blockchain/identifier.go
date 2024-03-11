package blockchain

import "github.com/google/uuid"

var NodeIdentifier string

func init() {
	NodeIdentifier = uuid.New().String()
}
