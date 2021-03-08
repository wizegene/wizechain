package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/davecgh/go-spew/spew"
	lru "github.com/karlseguin/ccache/v2"
)

type State struct {
	Pk             int `storm:"id,increment"`
	IsGenesisState bool
	Balances       map[string]float32
	txMemPool      []Transaction
	cache          *lru.Cache
	IState
}

func (s *State) CreateCache() {

	cache := lru.New(lru.Configure().MaxSize(10000).ItemsToPrune(100))
	s.cache = cache

}

func (s *State) AddToMempool(key string, value interface{}) {
	s.cache.Set(key, value, time.Hour*1)
}

func (s *State) GetCacheSize() int {
	return s.cache.ItemCount()
}

type IState interface {
	GetDatabase() *storm.DB
}

var stateDB *storm.DB

func (s *State) GetDatabase() *storm.DB {

	stateDB, _ := storm.Open("./db/states.db")
	return stateDB

}

func (s *State) CreateGenesisState(db *storm.DB) {
	defer db.Close()

	err := db.Init(&State{})
	if err != nil {
		fmt.Printf("error: %v", err)

	}

	var state State
	err = db.One("IsGenesisState", true, &state)
	if err == storm.ErrNotFound {

		//ok genesis state does not exists we read the json genesis file
		g, _ := ioutil.ReadFile("./config/genesis.json")
		var G Genesis
		err2 := json.Unmarshal(g, &G)
		if err2 != nil {
			fmt.Printf("error: %v", err2)
			panic(err2)
		}

		state.Balances = G.Balances
		state.IsGenesisState = true
		err2 = db.Save(&state)
		if err2 != nil {
			panic(err2)
		}

	} else {
		// genesis state exists we do not create it
		var states []State
		err := db.All(&states)
		if err != nil {
			panic(err)
		}
		spew.Dump(states)

	}
}

func init() {
	s := new(State)
	db := s.GetDatabase()
	s.CreateGenesisState(db)
}

//AddTx ...
func (s *State) AddTx(tx Transaction) error {

	s.txMemPool = append(s.txMemPool, tx)

	return nil

}

//Persist ...
func (s *State) Persist() error {

	var tx Transaction
	mempool := make([]Transaction, len(s.txMemPool))
	copy(mempool, s.txMemPool)
	db := tx.GetDatabase()
	defer db.Close()

	for i := 0; i < len(mempool); i++ {
		if err = db.Save(s.txMemPool[i]); err != nil {
			return err
		}

		s.txMemPool = append(s.txMemPool[:i], s.txMemPool[i+1:]...)

	}

	return nil
}

func (s *State) apply(tx Transaction) error {
	if tx.IsStakingReward {
		s.Balances[tx.Inputs[0].To] += tx.Inputs[0].Value
		return nil
	}

	if tx.Inputs[0].Value > s.Balances[tx.Inputs[0].From] {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.Inputs[0].From] -= tx.Inputs[0].Value
	s.Balances[tx.Inputs[0].To] += tx.Inputs[0].Value
	return nil
}

func (s *State) getAllBalances() (*State, error) {
	var tx Transaction

	var txs []Transaction
	db := tx.GetDatabase()
	defer db.Close()
	db.All(&txs)
	state := &State{Balances: make(map[string]float32, 0)}

	for _, tx := range txs {

		if err := state.apply(tx); err != nil {
			return nil, err
		}

	}
	return state, nil

}
