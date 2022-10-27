package models

type User struct {
	Login     string
	Password  string
	Email     string
	ProfileId uint32
	FirstName string
	LastName  string
}

type Profile struct {
	ProfileId uint32
	AvatarUrl string
	FirstName string
	LastName  string
}
