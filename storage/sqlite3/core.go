package sqlite3

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func NewDB(databaseName string) (*sql.DB, error) {

	if _, err := os.Stat("/path/to/whatever"); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(databaseName)
		if err != nil {
			return nil, err
		}
		err = file.Close()
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return nil, err
	}

	_, err = MigrateToActualVersion(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		return
	}
}
