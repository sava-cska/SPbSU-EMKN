package models

type User struct {
	Login     string
	Password  string
	Email     string
	ProfileId int32
	FirstName string
	LastName  string
}

type Profile struct {
	ProfileId int32
	AvatarUrl string
	FirstName string
	LastName  string
}
