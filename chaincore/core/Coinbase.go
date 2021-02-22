package core

import "sync"

type ICoinbase interface {
	clear()
	setDefaults()
	Cleanup()
	ClearUnspendable()
	CalcMaskSize()
	Serialize() []byte
	UnSerialize() *Coin
	isSpent(pos uint32) bool
	markAsSpent(pos uint32)
	isAvailable(pos uint32) bool
	isPruned() bool
	GetAddress() string
	GetUTXO(txID string) []Coin
	HaveUTXO(txID string) bool
	GetCoinBlock() *Block
	batchWrite(coinMap map[uint]Coin, blockHash []byte)
	SetCacheBackend()
	haveCoinsInCache(txID string) bool
	accessCoinsFromCache(txID string) []CoinCacheEntry
	applyInternalFeePolicy()
	_addCoinbaseTX()
	_getCoinbaseTX(txID string)
	_mint()
	_burn()
	_transfer(from string, to string, value []byte)
	_swap(coins []Coin, to []Coin)
	_lock()
	_unlock()
	_pause()
	_unpause()
	_modifyCoins(txID string)
	_modifyNewCoins(txID string, isCoinbase bool)
	_flush() bool
	_unCache(txID string)
	_getCacheSize() int32
	_getCacheSizeInBytes() uint32
	_getValueIn(tx *Transaction) float64
	_haveInputs(tx *Transaction) bool
	_getPriority(tx *Transaction, height uint32)
	_getOutputFor(tx *Transaction) Outputs
	_fetchCoins(txID string) map[uint32]Coin
}

type Coinbase struct {
	address      string
	cacheBackend *CoinCache
	coins        []*Coin
	mu           sync.Mutex
}

type CoinCache struct {
	mu             sync.Mutex
	Entries        map[uint32]CoinCacheEntry
	hasModifier    bool
	cacheCoinUsage uint32
}

type CoinCacheEntry struct {
	coin      Coin
	cacheFlag bool
}

const (
	CoinFlagDirty = 0
	CoinFlagFresh = 1
)

type CoinbaseManager struct {
	ICoinbase
	coinbase *Coinbase
}

func NewCoinbase() *CoinbaseManager {
	return &CoinbaseManager{}
}

func (cb *CoinbaseManager) CalcMaskSize() {

}

func (cb *CoinbaseManager) SetCacheBackend() {

	cCacheView := &CoinCache{}
	cCacheView.mu.Lock()
	cCacheView.hasModifier = false
	cCacheView.Entries = make(map[uint32]CoinCacheEntry, 0)
	cCacheView.cacheCoinUsage = 0
	cb.coinbase.cacheBackend = cCacheView

}

func (cb *CoinbaseManager) haveCoinsInCache(txID string) bool {
	return false
}

func (cb *CoinbaseManager) isSpent(pos uint32) bool {
	return true
}

func (cb *CoinbaseManager) setDefaults() {

	cb.coinbase.address = Encode([]byte("0xfffffffffffffffffffffffffffffff"), WizegeneAlphabet)

}

func (cb *CoinbaseManager) clear() {

}
