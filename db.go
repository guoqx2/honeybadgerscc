package main

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

func dbPut(key string, value string) {
	db, err := leveldb.OpenFile("db", nil)
	err = db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		fmt.Println("Error writing to database")
	}
	defer db.Close()
}

func dbGet(key string) string {
	db, err := leveldb.OpenFile("db", nil)
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		fmt.Println("Error getting from database")
	}
	defer db.Close()
	return string(data)
}
