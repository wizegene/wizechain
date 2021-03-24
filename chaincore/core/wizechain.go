package core

import (
	"encoding/hex"
	"encoding/json"
	db "github.com/wizegene/wizechain/chaincore/database"
	"io/ioutil"
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
	bQueue            *BlockQueue
	MasterFingerPrint string
	MasterDNA         string
	memPool           map[string][]byte
}

type IWizechain interface {
	authorize()
	initNewChain(id string, chaincode string, version string) *Wizechain
	GetGenesis()
	SetGenesis()
	CreateBlock()
	GetTip() *Block
	GetEnd() *Block
	GetNumBlocks() int
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

func NewWizeChain(chaincode, version string) *Wizechain {

	var ws Wizechain
	return ws.initNewChain(chaincode, version)
}

func (w *Wizechain) initNewChain(chaincode string, version string) *Wizechain {
	ChainDB := db.InitDB(chaincode + "_" + version + ".db")

	defer ChainDB.Close()
	w.memPool = make(map[string][]byte, 0)
	masterSeed, err = CreateSeedForKey()
	w.memPool["master_seed"] = masterSeed
	if err != nil {
		panic(err)
	}

	w.SetGenesis()
	w.bQueue = NewBlockQueue()

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

	return w

}

func (w *Wizechain) GetGenesis() {
	LoadGenesisFromJSON()
}

func (w *Wizechain) SetGenesis() {
	var block Block
	ser, _ := ioutil.ReadFile("./config/in_block.json")
	_ = json.Unmarshal(ser, &block)
	w.Blocks = make([]*Block, 0)
	w.Blocks = append(w.Blocks, &block)
	w.ID = block.Header.ChainID
}

func (w *Wizechain) CreateBlock() {

	lastBlock := w.GetPreviousBlock()
	block := CreateRawBlock(w.ID, []byte(lastBlock.Header.Hash), lastBlock.Header.Id, lastBlock.Header.Version)
	w.bQueue.Add(block)
}

func (w *Wizechain) GetTip() *Block {
	return w.Blocks[0]
}

func (w *Wizechain) GetEnd() *Block {
	return w.Blocks[len(w.Blocks)-1]
}

func (w *Wizechain) GetNumBlocks() int {
	return 0
}

func getChain(id string, chaincode string, version string) {

}

func (w *Wizechain) GetPreviousBlock() *Block {
	return w.Blocks[w.GetNumBlocks()-1]
}
