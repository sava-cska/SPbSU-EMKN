package storage

import (
	"time"
)

type RegistrationDao struct {
	Storage *Storage
}

func (dao *RegistrationDao) Upsert(
	token string,
	login string,
	password string,
	email string,
	firstName string,
	lastName string,
	expireDate time.Time,
	verificationCode string,
) error {
	_, err := dao.Storage.Db.Exec(
		"INSERT INTO\n"+
			"registration_base (token, login, password, email, first_name, last_name, expire_date, verification_code)\n"+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		token,
		login,
		password,
		email,
		firstName,
		lastName,
		expireDate,
		verificationCode,
	)
	return err
}
