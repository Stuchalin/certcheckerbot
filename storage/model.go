package storage

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
