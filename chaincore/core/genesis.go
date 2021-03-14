package core

import (
	"encoding/json"
	"errors"
	"github.com/beevik/ntp"
	"github.com/rogpeppe/fastuuid"
	"github.com/wizegene/wizechain/chaincore/config"
	"io/ioutil"
)

type Genesis struct {
	GenesisTime      int64              `json:"genesis_time"`
	ChainID          string             `json:"chain_id"`
	Balances         map[string]float32 `json:"balances"`
	Currency         *CurrencyConfig    `json:"currency"`
	POS              *POSParams         `json:"pos_parameters"`
	BlockchainParams *BlocksParams      `json:"blocks_parameters"`
	inBlock          []byte
}

type CurrencyConfig struct {
	Name      string
	MaxSupply int `json:"max_supply"`
}

func LoadGenesisFromJSON() *Genesis {

	var G Genesis
	g, _ := ioutil.ReadFile("./config/genesis.json")

	_ = json.Unmarshal(g, &G)
	return &G

}

func WriteGenesisToFile(data *Genesis) {
	ser, _ := json.Marshal(data)
	_ = ioutil.WriteFile("./config/genesis.json", ser, 0644)

}

func WriteInBlockToFile(blockData []byte) {
	_ = ioutil.WriteFile("./config/in_block.json", blockData, 0644)

}

func CreateGenesisBlock() (*Genesis, error) {

	g := &Genesis{}

	chainID, err := fastuuid.NewGenerator()
	if err != nil {
		return nil, err
	}
	// new randomize uuid 192 bit as hex 128 bit string
	cid := chainID.Hex128()
	//we create a blank block to be associated to the genesis content

	bser := CreateNewBlock(cid, nil, uint32(0), uint32(1))
	var block Block
	json.Unmarshal(bser, &block)

	blockTime := block.Created
	timestamp, _ := ntp.Time("pool.ntp.org")
	if blockTime > timestamp.UnixNano() {
		//invalid time
		return nil, errors.New("invalid genesis time")
	}
	if blockTime == timestamp.UnixNano() {
		//impossible blocktime for the block cannot be equal to genesis time because genesis is created after
		return nil, errors.New("invalid genesis time")
	}

	g.GenesisTime = timestamp.UnixNano()
	g.ChainID = cid
	g.Balances = make(map[string]float32)
	g.Currency = &CurrencyConfig{}
	g.POS = &POSParams{}
	g.BlockchainParams = &BlocksParams{}
	g.inBlock = bser

	WriteInBlockToFile(bser)

	g.Currency.Name = config.Currency
	g.Currency.MaxSupply = -1
	g.BlockchainParams.CoinsPerBlock = int32(config.CoinsPerBlock)
	g.BlockchainParams.isInflationary = true
	g.BlockchainParams.MaxBlockPerYear = float32(config.MaxBlockPerYear)
	g.BlockchainParams.MaxSecPerBlock = int64(config.SecondsBetweenBlocks)
	g.POS = &POSParams{}
	g.POS.MaxValidators = config.MaxValidators
	g.POS.MinValidators = config.MinValidators
	g.POS.InflationaryRate = config.RewardInflationPerSecond
	g.POS.StartingRewardPerValidator = config.StartRewardPerValidator
	g.POS.StartingRewardAsCurrency = config.StartRewardAsWize
	g.POS.Quorum = float32(config.Corum)

	WriteGenesisToFile(g)

	return g, nil

}

func (g *Genesis) toJSON() []byte {

	ser, _ := json.Marshal(g)
	return ser

}

type POSParams struct {
	MinValidators              int
	MaxValidators              int
	Quorum                     float32
	MaxValidatorPower          float64
	ChangeLeaderFreq           float64
	StartingRewardPerValidator float64
	StartingRewardAsCurrency   float64
	InflationaryRate           float64
	DeflationaryRate           float64
}

type BlocksParams struct {
	MaxBlockPerYear float32
	CoinsPerBlock   int32
	MaxCoinsPerYear float64
	isDeflationary  bool
	isInflationary  bool
	MaxSecPerBlock  int64
}
