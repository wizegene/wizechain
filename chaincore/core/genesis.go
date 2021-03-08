package core

import (
	"encoding/json"
	"io/ioutil"
)

type Genesis struct {
	GenesisTime string             `json:"genesis_time"`
	ChainID     string             `json:"chain_id"`
	Balances    map[string]float32 `json:"balances"`
	Currency    *CurrencyConfig    `json:"currency"`
}

type CurrencyConfig struct {
	Name      string
	MaxSupply int `json:"max_supply"`
}

func LoadGenesisFromJSON() *Genesis {

	var G Genesis
	g, _ := ioutil.ReadFile("./config/genesis.json")

	json.Unmarshal(g, &G)
	return &G

}
