package database

import (
	"chaincore/config"
	badger "github.com/dgraph-io/badger/v3"
)

var db *badger.DB
var err error

func init() {
	db, err = badger.Open(badger.DefaultOptions(config.BadgerDBDir))
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

type Database struct {
	db *badger.DB
}

func Insert(key []byte, value []byte) error {
	defer db.Close()
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
	if err != nil {
		return err
	}
	return nil
}

func FetchByKey(key []byte) ([]byte, error) {
	defer db.Close()

	var itemString string
	_ = db.View(func(txn *badger.Txn) error {
		Item, err := txn.Get(key)
		if err != nil {
			return err
		}
		itemString = Item.String()
		return nil

	})
	if err != nil {
		return nil, err
	}
	return []byte(itemString), nil
}
