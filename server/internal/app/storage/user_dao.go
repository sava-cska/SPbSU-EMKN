package storage

import (
	"database/sql"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
)

type UserDAO struct {
	Storage *Storage
}

func (dao *UserDAO) ExistsLogin(login string) bool {
	row := dao.Storage.Db.QueryRow(
		`SELECT login
			   FROM user_base
			   WHERE login = $1`,
		login,
	)
	var tmpLogin string
	err := row.Scan(&tmpLogin)
	return err != sql.ErrNoRows
}

func (dao *UserDAO) ExistsEmail(email string) bool {
	row := dao.Storage.Db.QueryRow(
		`SELECT email
			   FROM user_base
			   WHERE email = $1`,
		email,
	)
	var tmpLogin string
	err := row.Scan(&tmpLogin)
	return err != sql.ErrNoRows
}

func (dao *UserDAO) AddUser(user *models.User) error {
	_, err := dao.Storage.Db.Exec(
		`INSERT INTO user_base (login, password, email, first_name, last_name)
               VALUES ($1, $2, $3, $4, $5)`,
		user.Login, user.Password, user.Email, user.FirstName, user.LastName)
	return err
}

func (dao *UserDAO) FindUser(email string) (models.User, error) {
	row := dao.Storage.Db.QueryRow(
		`
		SELECT login, password, email, first_name, last_name
		FROM user_base
		WHERE email = $1
		`,
		email)

	var user models.User
	err := row.Scan(&user.Login, &user.Password, &user.Email, &user.FirstName, &user.LastName)
	return user, err
}

func (dao *UserDAO) FindUserByLogin(login string) (models.User, error) {
	row := dao.Storage.Db.QueryRow(
		`
		SELECT login, password, email, first_name, last_name
		FROM user_base
		WHERE login = $1
		`,
		login)

	var user models.User
	err := row.Scan(&user.Login, &user.Password, &user.Email, &user.FirstName, &user.LastName)
	return user, err
}

func (dao *UserDAO) GetPassword(login string) (string, error) {
	row := dao.Storage.Db.QueryRow(
		`SELECT password
               FROM user_base
               WHERE login = $1
    `, login)

	var pwd string
	err := row.Scan(&pwd)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return pwd, err
}

func (dao *UserDAO) UpdatePassword(login string, newPassword string) error {
	_, errUpdate := dao.Storage.Db.Exec(
		`
			UPDATE user_base
			SET
			password = $1
			WHERE login = $2
		`,
		newPassword, login)
	return errUpdate
}
