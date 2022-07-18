package sqlite3

import (
	"certcheckerbot/storage"
	"database/sql"
	"fmt"
)

type Sqlite3Controller struct {
	Connection *sql.DB
}

func NewController(databaseName string) (*Sqlite3Controller, error) {
	db, err := NewDB(databaseName)
	if err != nil {
		return nil, err
	}
	result := Sqlite3Controller{
		Connection: db,
	}
	return &result, nil
}

func (db *Sqlite3Controller) AddUser(user *storage.User) (int, error) {
	stmt, err := db.Connection.Prepare("insert into Users(Name, TGId, NotificationHour, UTC) values (?, ?, ?, ?);")
	if err != nil {
		return 0, err
	}
	tx, err := db.Connection.Begin()
	if err != nil {
		return 0, err
	}
	result, err := tx.Stmt(stmt).Exec(user.Name, user.TGId, user.NotificationHour, user.UTC)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	user.Id = int(id)

	return int(id), nil
}

func (db *Sqlite3Controller) GetUserById(id int) (*storage.User, error) {
	record, err := db.Connection.Query("select Id, Name, TGId, NotificationHour, UTC from Users where Id = ?;", id)
	if err != nil {
		return nil, err
	}
	defer func(record *sql.Rows) {
		_ = record.Close()
	}(record)

	if record.Next() {
		var user storage.User
		err := record.Scan(&user.Id, &user.Name, &user.TGId, &user.NotificationHour, &user.UTC)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}

	return nil, fmt.Errorf("cannot find user by id %d", id)
}

func DisposeConnection(db *Sqlite3Controller) {
	CloseConnection(db.Connection)
}
