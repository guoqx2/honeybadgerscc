package main

import (
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var mut sync.Mutex

func dbPut(db *leveldb.DB, key string, value string) {
	mut.Lock()
	if db == nil {
		fmt.Println("db is nil ")
		db, _ = leveldb.OpenFile("db", nil)
	}
	err := db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		fmt.Println("Error writing to database")
	}
	db.Close()
	mut.Unlock()

}

func dbGet(db *leveldb.DB, key string) string {
	mut.Lock()
	if db == nil {
		fmt.Println("db is nil")
		db, _ = leveldb.OpenFile("db", nil)
	}
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		fmt.Println("Error getting from database")
	}
	db.Close()
	mut.Unlock()

	return string(data)
}
