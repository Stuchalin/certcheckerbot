package sqlite3

import (
	"database/sql"
	"errors"
)

type Migration struct {
	Version         int
	MigrationScript string
}

func getMigrations() []Migration {
	return []Migration{
		{Version: 1, MigrationScript: "" +
			"CREATE TABLE db_versions (" +
			"	Version INTEGER," +
			"	PRIMARY KEY (Version)" +
			");" +
			"CREATE TABLE Users (" +
			"	Id INTEGER," +
			"	Name VARCHAR(4000)," +
			"	TGId VARCHAR(255)," +
			"	NotificationHour INTEGER," +
			"	UTC INTEGER," +
			"	PRIMARY KEY (Id)" +
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

type User struct {
	Id               int
	Name             string
	TGId             string
	NotificationHour int
	UTC              int
	UserDomains      []UserDomain
}

type UserDomain struct {
	UserId int
	Domain string
}

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
		var version int
		record.Scan(&version)
		return version, nil
	} else {
		return 0, nil
	}
}

//func migrateDatabase(migration Migration, db *sql.DB) (bool, error) {
//	tx, err := db.Begin()
//	if err != nil {
//		return false, err
//	}
//
//	tx.Exec(migration)
//}

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
			tx.Rollback()
		}
		return false, err
	}
	_, err = stmt.Exec(version)
	if err != nil {
		if runInPersonalTransaction {
			tx.Rollback()
		}
		return false, err
	}
	if runInPersonalTransaction {
		tx.Commit()
	}
	return true, nil
}
