package main

import (
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	db, err := leveldb.OpenFile("test_db", nil)
	check(err)
	defer db.Close()

	err = db.Put([]byte("1"), []byte("john"), nil)
	check(err)
	err = db.Put([]byte("2"), []byte("joe"), nil)
	check(err)
	err = db.Put([]byte("3"), []byte("jane"), nil)
	check(err)
	err = db.Put([]byte("4"), []byte("janet"), nil)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
