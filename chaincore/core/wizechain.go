package core

/*
This is the main file to generate the wizechain blockchain
*/

import "C"
import (
	db "chaincore/database"
	"time"
)

var masterSeed []byte
var err error

func initNewChain(id string, chaincode string, version string) {
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

func getChain(id string, chaincode string, version string) {

}
