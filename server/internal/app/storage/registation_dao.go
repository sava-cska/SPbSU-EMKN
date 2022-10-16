package storage

import (
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
)

type RegistrationDAO struct {
	Storage *Storage
}

func (dao *RegistrationDAO) Upsert(
	token string,
	user *models.User,
	expireDate time.Time,
	verificationCode string,
) error {
	_, err := dao.Storage.Db.Exec(
		`INSERT INTO
			   registration_base (token, login, password, email, first_name, last_name, expire_date, verification_code)
		       VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		token,
		user.Login,
		user.Password,
		user.Email,
		user.FirstName,
		user.LastName,
		expireDate,
		verificationCode,
	)
	return err
}

func (dao *RegistrationDAO) FindRegistration(token string) (models.User, time.Time, string, error) {
	registerRecord := dao.Storage.Db.QueryRow(
		`SELECT login, password, email, first_name, last_name, expire_date, verification_code
			   FROM registration_base WHERE token = $1`,
		token)

	var user models.User
	var expireTime time.Time
	var verificationCode string
	if errScan := registerRecord.Scan(&user.Login, &user.Password, &user.Email, &user.FirstName,
		&user.LastName, &expireTime, &verificationCode); errScan != nil {
		return models.User{}, time.Time{}, "", errScan
	}
	return user, expireTime, verificationCode, nil
}

func (dao *RegistrationDAO) FindRegistrationAndDelete(token string) (models.User, time.Time, string, error) {
	tx, err := dao.Storage.Db.Begin()
	if err != nil {
		return models.User{}, time.Time{}, "", err
	}
	registerRecord := tx.QueryRow(
		`SELECT login, password, email, first_name, last_name, expire_date, verification_code
			   FROM registration_base WHERE token = $1`,
		token)

	var user models.User
	var expireTime time.Time
	var verificationCode string
	if errScan := registerRecord.Scan(&user.Login, &user.Password, &user.Email, &user.FirstName,
		&user.LastName, &expireTime, &verificationCode); errScan != nil {
		tx.Rollback()
		return models.User{}, time.Time{}, "", errScan
	}
	_, err = tx.Exec(
		`DELETE FROM registration_base WHERE token = $1`,
		token,
	)
	if err != nil {
		tx.Rollback()
		return models.User{}, time.Time{}, "", err
	}
	tx.Commit()
	return user, expireTime, verificationCode, nil
}
