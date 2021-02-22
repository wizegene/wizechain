package core

import "sync"

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
