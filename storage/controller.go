package storage

import "errors"

var ErrorUserNotFound = errors.New("storage error - user not found")
var ErrorUserDomainNotFound = errors.New("storage error - user domain not found")

type UsersConfig interface {
	AddUser(user *User) (int, error)
	GetUserById(id int) (*User, error)
	GetUserByTGId(tgId int64) (*User, error)
	GetUserByName(name string) (*User, error)
	RemoveUser(user *User) (bool, error)
	UpdateUserInfo(user *User) (bool, error)

	AddUserDomain(domain *UserDomain) (bool, error)
	RemoveUserDomain(domain *UserDomain) (bool, error)
	GetUserDomains(user *User) (*[]UserDomain, error)
}
