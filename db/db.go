package db

import (
	"log"

	"github.com/dgraph-io/badger/v3"
)

var database *badger.DB

func InitDB() error {
	opts := badger.DefaultOptions("./data")
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	database = db
	return nil
}

func GetDB() *badger.DB {
	return database
}

func CloseDB() {
	if database != nil {
		err := database.Close()
		if err != nil {
			log.Println("Error closing database:", err)
		}
	}
}
