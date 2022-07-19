package sqlite3

import (
	"certcheckerbot/storage"
	"database/sql"
)

//Sqlite3Controller controller for sqlite3 database
type Sqlite3Controller struct {
	Connection *sql.DB
}

//NewController creates new database controller with connection
//databasePath - path to database
func NewController(databasePath string) (*Sqlite3Controller, error) {
	db, err := NewDB(databasePath)
	if err != nil {
		return nil, err
	}
	result := Sqlite3Controller{
		Connection: db,
	}
	return &result, nil
}

//AddUser add new user in database
func (db *Sqlite3Controller) AddUser(user *storage.User) (int, error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return -1, err
	}
	userId, err := addUser(user, tx)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	user.Id = userId

	return userId, nil
}

//addUser - add user processing, expected external transaction
func addUser(user *storage.User, tx *sql.Tx) (int, error) {
	stmt, err := tx.Prepare("insert into Users(Name, TGId, NotificationHour, UTC) values (?, ?, ?, ?);")
	if err != nil {
		return -1, err
	}
	result, err := tx.Stmt(stmt).Exec(user.Name, user.TGId, user.NotificationHour, user.UTC)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

//AddUserDomain - add tracked domain to user
func (db *Sqlite3Controller) AddUserDomain(domain *storage.UserDomain) (bool, error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return false, err
	}
	result, err := addUserDomain(domain, tx)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return result, nil
}

//addUserDomain - add tracked domain to user processing, expected external transaction
func addUserDomain(domain *storage.UserDomain, tx *sql.Tx) (bool, error) {
	stmt, err := tx.Prepare("insert into UserDomains(UserId, Domain) values (?, ?);")
	if err != nil {
		return false, err
	}
	result, err := tx.Stmt(stmt).Exec(domain.UserId, domain.Domain)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected > 0 {
		return true, nil
	}
	return false, nil
}

//RemoveUserDomain - remove tracked domain from user
func (db *Sqlite3Controller) RemoveUserDomain(domain *storage.UserDomain) (bool, error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return false, err
	}
	result, err := removeUserDomain(domain, tx)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return result, nil
}

//removeUserDomain - remove tracked domain from user, expected external transaction
func removeUserDomain(domain *storage.UserDomain, tx *sql.Tx) (bool, error) {
	stmt, err := tx.Prepare("delete from UserDomains where UserId = ? and Domain = ?;")
	if err != nil {
		return false, err
	}
	result, err := tx.Stmt(stmt).Exec(domain.UserId, domain.Domain)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected > 0 {
		return true, nil
	}
	return false, nil
}

//RemoveUser - remove user from database
func (db *Sqlite3Controller) RemoveUser(user *storage.User) (bool, error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return false, err
	}
	_, err = removeAllUserDomains(user, tx)
	result, err := removeUser(user, tx)
	if err != nil {
		return false, err
	}
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return result, nil
}

//removeUser - remove user from database processing, expected external transaction
func removeUser(user *storage.User, tx *sql.Tx) (bool, error) {
	stmt, err := tx.Prepare("delete from Users where Id = ?;")
	if err != nil {
		return false, err
	}
	result, err := tx.Stmt(stmt).Exec(user.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected > 0 {
		return true, nil
	}
	return false, nil
}

//removeAllUserDomains - remove all user domains from database processing, expected external transaction
func removeAllUserDomains(user *storage.User, tx *sql.Tx) (bool, error) {
	stmt, err := tx.Prepare("delete from UserDomains where UserId = ?;")
	if err != nil {
		return false, err
	}
	result, err := tx.Stmt(stmt).Exec(user.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected > 0 {
		return true, nil
	}
	return false, nil
}

//UpdateUserInfo - updates user info in database
func (db *Sqlite3Controller) UpdateUserInfo(user *storage.User) (bool, error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return false, err
	}
	result, err := updateUserInfo(user, tx)
	if err != nil {
		return false, err
	}
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return result, nil
}

//updateUserInfo - updates user info in database processing, expected external transaction
func updateUserInfo(user *storage.User, tx *sql.Tx) (bool, error) {
	stmt, err := tx.Prepare("update Users" +
		"	set Name = ?," +
		"	TGId = ?," +
		"	NotificationHour = ?," +
		"	UTC = ?" +
		"where Id = ?;")
	if err != nil {
		return false, err
	}
	result, err := tx.Stmt(stmt).Exec(user.Name, user.TGId, user.NotificationHour, user.UTC, user.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected > 0 {
		return true, nil
	}
	return false, nil
}

//GetUserById - search user from database by id
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

	return nil, storage.ErrorUserNotFound
}

//GetUserByTGId - search user from database by TGId
func (db *Sqlite3Controller) GetUserByTGId(tgId int64) (*storage.User, error) {
	record, err := db.Connection.Query("select Id, Name, TGId, NotificationHour, UTC from Users where TGId = ?;", tgId)
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

	return nil, storage.ErrorUserNotFound
}

//GetUserByName - search user from database by name
func (db *Sqlite3Controller) GetUserByName(name string) (*storage.User, error) {
	record, err := db.Connection.Query("select Id, Name, TGId, NotificationHour, UTC from Users where Name = ?;", name)
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

	return nil, storage.ErrorUserNotFound
}

//GetUserDomains - select user domains from database
func (db *Sqlite3Controller) GetUserDomains(user *storage.User) (*[]storage.UserDomain, error) {
	record, err := db.Connection.Query("select UserId, Domain from UserDomains where UserId = ?;", user.Id)
	if err != nil {
		return nil, err
	}
	defer func(record *sql.Rows) {
		_ = record.Close()
	}(record)

	var userDomains []storage.UserDomain

	for record.Next() {
		var userDomain storage.UserDomain
		err := record.Scan(&userDomain.UserId, &userDomain.Domain)
		if err != nil {
			return nil, err
		}
		userDomains = append(userDomains, userDomain)
	}
	if userDomains != nil {
		return &userDomains, nil
	}

	return nil, storage.ErrorUserDomainNotFound
}

//Dispose - close connections to database
func (db *Sqlite3Controller) Dispose() {
	CloseConnection(db.Connection)
}
