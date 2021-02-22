package core

import (
	"sync"
)

type ITransaction interface {
	isCoinbase() bool
	fromTx(txID string, height uint32)
}

type Transaction struct {
	ITransaction
	mu sync.Mutex
}

type Outputs struct {
}

type Inputs struct {
}

type Outpoint struct {
	Hash []byte
	n    uint32
}

// NewOutpoint ...
func NewOutpoint(hashIn []byte, nIn uint32) *Outpoint {
	return &Outpoint{
		Hash: hashIn,
		n:    nIn,
	}
}

func (o *Outpoint) Serialize() {

}
