package sqlite3

import (
	"database/sql"
	"errors"
)

func isBaseStructExists(db *sql.DB) (bool, error) {
	if db == nil {
		return false, errors.New("Nil database connection")
	}

	record, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='db_version';")
	if err != nil {
		return false, err
	}
	defer record.Close()

	if record.Next() {
		return true, nil
	} else {
		return false, nil
	}
}

func GetCurrentDBVersion(db *sql.DB) (int, error) {

	baseStruct, err := isBaseStructExists(db)

	if err != nil {
		return -1, err
	}

	if !baseStruct {
		return 0, nil
	} else {

	}

	record, err := db.Query("SELECT max(version) FROM db_version;")
	if err != nil {
		return -1, err
	}
	defer record.Close()

	if record.Next() {
		var version int
		record.Scan(&version)
		return version, nil
	} else {
		return 0, nil
	}
}
