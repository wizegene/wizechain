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
	memPool      map[uint32][]byte
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
	return &Wizechain{}
}

func (w *Wizechain) initNewChain(id string, chaincode string, version string) {
	ChainDB := db.InitDB(id + "/" + chaincode + "_" + version)

	defer ChainDB.Close()
	db.Insert([]byte("_chain__initialization_time"), []byte(time.Now().String()))
	masterSeed, err = CreateSeedForKey()
	if err != nil {
		panic(err)
	}

	masterKey := CreateMasterKey(masterSeed)
	ms := masterKey.B58Serialize()
	db.Insert([]byte("_chain__master_key_genesis"), []byte(ms))
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
