package core

import (
	"encoding/hex"
	db "github.com/wizegene/wizechain/chaincore/database"
	"sync"
)

var masterSeed []byte
var err error

type Wizechain struct {
	mu                *sync.Mutex
	ID                string
	Chaincode         string
	Version           string
	Blocks            []*Block
	BlockHeaders      []*BlockHeader
	MasterFingerPrint string
	MasterDNA         string
	memPool           map[string][]byte
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
	ChainDB := db.InitDB(id + "/" + chaincode + "_" + version + ".db")

	defer ChainDB.Close()

	masterSeed, err = CreateSeedForKey()
	w.memPool["master_seed"] = masterSeed
	if err != nil {
		panic(err)
	}

	masterKey := CreateMasterKey(masterSeed)
	w.memPool["master_key"] = []byte(masterKey.B58Serialize())

	ms := masterKey.FingerPrint

	dna := GetDNA(5000)

	w.MasterDNA = hex.EncodeToString(dna.D)
	w.MasterFingerPrint = hex.EncodeToString(ms)
	err := ChainDB.Save(w)
	if err != nil {
		panic(err)
	}

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
