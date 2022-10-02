package storage

import (
	"database/sql"
)

type UserDAO struct {
	Storage *Storage
}

func (dao *UserDAO) Exists(login string) bool {
	row := dao.Storage.Db.QueryRow(
		"SELECT\n"+
			"	login\n"+
			"FROM user_base\n"+
			"WHERE login = $1\n",
		login,
	)
	var tmpLogin string
	err := row.Scan(&tmpLogin)
	return err != sql.ErrNoRows
}
