package core

import (
	"errors"

	"github.com/asdine/storm/v3"
)

type Account Address

type ITransaction interface {
	isCoinbase() bool
	fromTx(txID string, height uint32)
}

type Transaction struct {
	Pk int `storm:"id,increment"`
	ITransaction
	Inputs          []*Inputs
	Outputs         []*Outputs
	IsCoinbase      bool
	IsStakingReward bool
}

type Outputs struct {
	From  string
	To    string
	Value float32
}

type Inputs struct {
	From  string  `json:"from"`
	To    string  `json:"to"`
	Value float32 `json:"value"`
	Data  string  `json:"data"`
}

type Outpoint struct {
	Hash []byte
	n    uint32
}

// NewOutpoint ...
func NewOutpoint(hashIn []byte, nIn uint32) *Outpoint {
	return &Outpoint{
		Hash: hashIn,
		n:    nIn,
	}
}

func (o *Outpoint) Serialize() {

}

func (t *Inputs) IsStakingReward() bool {
	return t.Data == "staking_reward"
}

func CreateGenesisTransaction() {

	g := LoadGenesisFromJSON()
	tx := new(Transaction)
	tx.IsCoinbase = true
	input := new(Inputs)
	input.From = Encode([]byte("0x000000000000000000000000000000"), WizegeneAlphabet)
	input.To = "DrkjsMvSTfVEAdmYUx/urTWkepn03u0gXeSvgRIPl2E="
	input.Value = g.Balances["DrkjsMvSTfVEAdmYUx/urTWkepn03u0gXeSvgRIPl2E="]
	input.Data = "genesis_transaction_0"

	output := new(Outputs)
	output.From = input.From
	output.To = input.To
	output.Value = input.Value

	tx.Inputs = make([]*Inputs, 1)
	tx.Inputs[0] = input
	tx.Outputs = make([]*Outputs, 1)
	tx.Outputs[0] = output

	db := tx.GetDatabase()
	defer db.Close()
	exists, _ := GetGenesisTransaction(db)

	//if genesis transaction does not exist create one
	if exists == false {
		err := db.Save(tx)
		if err != nil {
			panic(err)
		}
	}

}

func (t *Transaction) GetDatabase() *storm.DB {

	stateDB, _ := storm.Open("./db/tx.db")
	return stateDB

}

func GetGenesisTransaction(db *storm.DB) (bool, []Transaction) {

	var txs []Transaction
	err := db.All(&txs)
	if err != nil {

		return false, nil
	}
	return true, txs

}

func init() {
	CreateGenesisTransaction()
}

func CreateFungibleTransaction(from string, to string, value float32) error {

	tx := new(Transaction)
	tx.IsCoinbase = false
	tx.Inputs = make([]*Inputs, 0)
	tx.Outputs = make([]*Outputs, 0)
	tx.Inputs[0].From = from
	tx.Inputs[0].To = to
	tx.Inputs[0].Value = value
	tx.Outputs[0] = new(Outputs)
	err := ValidateTXBeforeCreate(tx)
	if err != nil {
		return err
	}

	return nil

}

func ValidateTXBeforeCreate(tx *Transaction) error {

	bal := GetBalanceFromAddress(tx.Inputs[0].From)
	if bal == 0 {
		return errors.New("balance is 0 cannot send tx")
	}
	if bal < tx.Inputs[0].Value {
		return errors.New("balance is not enough to cover tx cannot send tx")
	}

	return nil

}

func GetBalanceFromAddress(from string) float32 {
	var state State
	db := state.GetDatabase()
	defer db.Close()
	var states []State
	err := db.All(&states)
	if err != nil {
		return 0
	}

	balance := float32(0)

	for i := 0; i < len(states); i++ {

		if states[i].Balances[from] != -1 {
			balance += states[i].Balances[from]
		}

	}

	return balance
}
