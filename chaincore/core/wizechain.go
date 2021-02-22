package core

import (
	db "chaincore/database"
	"sync"
	"time"
)

var masterSeed []byte
var err error

type Wizechain struct {
	mu           *sync.Mutex
	ID           string
	Chaincode    string
	Version      string
	Blocks       []*Block
	BlockHeaders []*BlockHeader
	memPool      map[string][]byte
	IWizechain
}

type IWizechain interface {
	authorize()
	initNewChain(id string, chaincode string, version string)
	GetGenesis()
	SetGenesis()
	CreateBlock() *Block
	GetTip()
	GetEnd()
	GetNumBlocks()
	getInitializationTime()
	validateInitializationTime() bool
	getMasterKey()
	getMasterDNA()
	GetBlockIndex(blockIndex uint32) Block
	GetBlockHeader(blockIndex uint32) Block
	GetChainHeight() uint32
	validate() bool
	sync()
	getValidators()
	getMedianTXTime()
}

func NewWizeChain() *Wizechain {
	mempool := make(map[string][]byte, 0)
	return &Wizechain{
		memPool: mempool,
	}
}

func (w *Wizechain) initNewChain(id string, chaincode string, version string) {
	ChainDB := db.InitDB(id + "/" + chaincode + "_" + version)

	defer ChainDB.Close()
	db.Insert([]byte("_chain__initialization_time"), []byte(time.Now().String()))
	masterSeed, err = CreateSeedForKey()
	w.memPool["master_seed"] = masterSeed
	if err != nil {
		panic(err)
	}

	masterKey := CreateMasterKey(masterSeed)
	w.memPool["master_key"] = []byte(masterKey.B58Serialize())

	ms := masterKey.FingerPrint
	db.Insert([]byte("_chain__master_key_genesis"), ms)
	dna := GetDNA(5000)
	db.Insert([]byte("_chain__master_dna_genesis"), dna.D)

}

func (w *Wizechain) GetGenesis() {

}

func (w *Wizechain) SetGenesis() {

}

func (w *Wizechain) CreateBlock() *Block {
	return &Block{}
}

func (w *Wizechain) GetTip() {

}

func (w *Wizechain) GetEnd() {

}

func getChain(id string, chaincode string, version string) {

}
