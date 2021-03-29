package core

import (
	"math/big"
	"time"
)

type WalletInterface interface {
}

type walletManager interface {
	Start()
	Close()
	CurrencyCode() string
	ExchangeRates()
	AddWatchedAddresses(addrs ...Address) error
	IsDust(amount big.Int) bool
	CurrentAddress(purpose int) Address
	NewAddress(purpose int) Address
	DecodeAddress(addr string) Address
	Balance() (confirmed, unconfirmed float64)
	Transactions() ([]Transaction, error)
	GetTransaction(txid string) (Transaction, error)
	ChainTip() (uint32, string)
	ReSyncBlockchain(fromTime time.Time)
	GetValidations(txid string) (validators, atHeight uint32, err error)
}

type walletKeyGen interface {
	ChildKey(keyBytes []byte, chaincode []byte, isPrivateKey bool)
	HasKey(addr Address)
}

type walletBanker interface {
	GetFeePerByte(feeLevel int) big.Int
	Spend(amount big.Int, addr Address, feeLevel int, referenceID string, spendAll bool)
}
