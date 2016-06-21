package api

import (
	"github.com/boltdb/bolt"
	"log"
	"path"
	"runtime"
	"time"
)

var dbName = "gost_internal.db"

// InternalDatabase is used for all internal storage requirements such as keeping track
// of FOI's by observation
type InternalDatabase struct {
	open bool
	bolt *bolt.DB
}

// Open opens the bolt database or creates if not exist
func (db *InternalDatabase) Open() error {
	var err error
	_, filename, _, _ := runtime.Caller(0) // get full path of this file
	dbFile := path.Join(path.Dir(filename), dbName)
	config := &bolt.Options{Timeout: 1 * time.Second}
	db.bolt, err = bolt.Open(dbFile, 0600, config)
	if err != nil {
		log.Fatal(err)
	}

	db.open = true

	return nil
}

func (db *InternalDatabase) Close() {
	db.open = false
	db.bolt.Close()
}
