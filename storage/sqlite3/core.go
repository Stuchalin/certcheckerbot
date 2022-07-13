package sqlite3

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func NewDB(databaseName string) (*sql.DB, error) {

	if _, err := os.Stat("/path/to/whatever"); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(databaseName)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}

func CloseConnection(db sql.DB) {
	db.Close()
}

//func migrateToActualVersion(db *sql.DB) (bool, error) {
//
//}
