package core

/*
This is the main file to generate the wizechain blockchain
 */

import (
	db "chaincore/database"
	"github.com/WizegeneFoundation/wizechain/dna-server/crypto"
	"time"
)

var masterSeed []byte
var err error

func init() {
	 db.Insert([]byte("_chain__initialization_time"), []byte(time.Now().String()))
	 masterSeed, err = CreateSeedForKey()
	 if err != nil {
	 	panic(err)
	 }

	 masterKey := CreateMasterKey(masterSeed)
	 ms := masterKey.B58Serialize()
	 db.Insert([]byte("_chain__master_key_genesis"), []byte(ms))
	 dna := crypto.GetDNA(5000)
}
