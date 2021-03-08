package database

import (
	"github.com/asdine/storm/v3"
)

var db *storm.DB
var err error

func InitDB(dbDir string) *storm.DB {
	db, err = storm.Open(dbDir)
	if err != nil {
		panic(err)
	}
	return db
}

type Database struct {
	db *storm.DB
}
