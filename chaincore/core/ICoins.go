package core

var (
	CoinsDefaultCoinbase = false
	CoinsDefaultOutput   = 0
	CoinsDefaultHeight   = uint32(0)
	CoinsDefaultVersion  = 0
)

type ICoins interface {
}

type Coin struct {
	ICoins
}
