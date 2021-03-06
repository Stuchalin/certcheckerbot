package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
)

//Migration type for migration
type Migration struct {
	Version         int
	MigrationScript string
}

//getMigrations slice of database migrations
//if you need change database struct - add new element to this array
//Migration.Version - must be a greater, then previous version Migration.Version++
//Migration.MigrationScript - must contain a valid migration script,
//that must migrate database from previous version, to the new specified one
//ADD ONLY SORTED ARRAY
func getMigrations() []Migration {
	return []Migration{
		{Version: 1, MigrationScript: "" +
			"CREATE TABLE db_versions (" +
			"	Version INTEGER," +
			"	PRIMARY KEY (Version)" +
			");" +
			"CREATE TABLE Users (" +
			"	Id INTEGER primary key autoincrement," +
			"	Name VARCHAR(4000)," +
			"	TGId INTEGER," +
			"	NotificationHour INTEGER," +
			"	UTC INTEGER" +
			");" +
			"CREATE UNIQUE INDEX IX_Users_TGId ON Users(TGId);" +
			"CREATE INDEX IX_Users_Name ON Users(Name);" +
			"CREATE TABLE UserDomains (" +
			"	UserId INTEGER," +
			"	Domain varchar(4000)," +
			"	PRIMARY KEY (UserId, Domain)," +
			"	FOREIGN KEY(UserId) REFERENCES Users(Id)" +
			");"},
	}
}

//isBaseStructExists - checks existing base struct for database
func isBaseStructExists(db *sql.DB) (bool, error) {
	if db == nil {
		return false, errors.New("Nil database connection")
	}

	record, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='db_versions';")
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

//GetCurrentDBVersion get version of database
func GetCurrentDBVersion(db *sql.DB) (int, error) {

	baseStruct, err := isBaseStructExists(db)

	if err != nil {
		return -1, err
	}

	if !baseStruct {
		return 0, nil
	}

	record, err := db.Query("SELECT max(Version) FROM db_versions;")
	if err != nil {
		return -1, err
	}
	defer record.Close()

	if record.Next() {
		var version sql.NullInt64
		err := record.Scan(&version)
		if err != nil {
			return -1, err
		}
		if !version.Valid {
			return 0, nil
		}
		return int(version.Int64), nil
	} else {
		return 0, nil
	}
}

//migrateDatabase migrate database to sent migration
func migrateDatabase(migration Migration, db *sql.DB) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	version, err := GetCurrentDBVersion(db)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}
	if version >= migration.Version {
		_ = tx.Rollback()
		return false, errors.New("database already migrate to this version")
	}

	_, err = tx.Exec(migration.MigrationScript)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}
	_, err = setDBVersion(migration.Version, db, tx)
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

//setDBVersion insert db version to database
//version - version to set
//db - database connection
//tx - go in transaction. [optional] if tx == nil then func open new tran and commit it when the work is completed
//else work will be in passed transaction
func setDBVersion(version int, db *sql.DB, tx *sql.Tx) (bool, error) {
	runInPersonalTransaction := false
	var err error
	if tx == nil {
		runInPersonalTransaction = true
		tx, err = db.Begin()
		if err != nil {
			return false, err
		}
	}
	stmt, err := tx.Prepare("insert into db_versions(version) values (?)")
	if err != nil {
		if runInPersonalTransaction {
			_ = tx.Rollback()
		}
		return false, err
	}
	_, err = stmt.Exec(version)
	if err != nil {
		if runInPersonalTransaction {
			_ = tx.Rollback()
		}
		return false, err
	}
	if runInPersonalTransaction {
		err := tx.Commit()
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

//MigrateToActualVersion migrate database to actual version
func MigrateToActualVersion(db *sql.DB) (bool, error) {
	version, err := GetCurrentDBVersion(db)
	if err != nil {
		return false, err
	}
	for _, migration := range getMigrations() {
		if version < migration.Version {
			_, err := migrateDatabase(migration, db)
			if err != nil {
				return false, fmt.Errorf("migration error - cannot migrate to version %d (%s)", migration.Version, err.Error())
			}
		}
	}
	return true, nil
}
