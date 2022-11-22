package storage

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
)

type UserDAO interface {
	ExistsLogin(login string) bool
	ExistsEmail(email string) bool
	AddUser(user *models.User) error
	FindUser(email string) (models.User, error)
	FindUserByLogin(login string) (models.User, error)
	GetPassword(login string) (string, error)
	UpdatePassword(login string, newPassword string) error
	FindUsers(ids []int32) ([]models.User, error)
}

type userDAO struct {
	Storage *DbStorage
}

func (dao *userDAO) ExistsLogin(login string) bool {
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

func (dao *userDAO) ExistsEmail(email string) bool {
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

func (dao *userDAO) AddUser(user *models.User) error {
	_, err := dao.Storage.Db.Exec(
		`INSERT INTO user_base (login, password, email, first_name, last_name)
               VALUES ($1, $2, $3, $4, $5)`,
		user.Login, user.Password, user.Email, user.FirstName, user.LastName)
	return err
}

func (dao *userDAO) FindUser(email string) (models.User, error) {
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

func (dao *userDAO) FindUserByLogin(login string) (models.User, error) {
	row := dao.Storage.Db.QueryRow(
		`
		SELECT login, password, email, profile_id, first_name, last_name
		FROM user_base
		WHERE login = $1
		`,
		login)

	var user models.User
	err := row.Scan(&user.Login, &user.Password, &user.Email, &user.ProfileId, &user.FirstName, &user.LastName)
	return user, err
}

func (dao *userDAO) GetPassword(login string) (string, error) {
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

func (dao *userDAO) UpdatePassword(login string, newPassword string) error {
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

func (dao *userDAO) FindUsers(ids []int32) ([]models.User, error) {
	var strProfileIds []string
	for _, val := range ids {
		strProfileIds = append(strProfileIds, strconv.Itoa(int(val)))
	}
	rows, err := dao.Storage.Db.Query(
		fmt.Sprintf(`SELECT profile_id, first_name, last_name
		FROM user_base
		WHERE profile_id IN (%s)`, strings.Join(strProfileIds, ", ")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ProfileId,
			&user.FirstName,
			&user.LastName); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}
