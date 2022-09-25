package storage

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

type UserDAO struct {
	Storage *Storage
}

func (dao *UserDAO) Exists(logger *logrus.Logger, login string) bool {
	logger.Debugf("%s", login)
	row := dao.Storage.Db.QueryRow(
		"SELECT\n"+
			"	login\n"+
			"FROM user_base\n"+
			"WHERE login = $1\n",
		login,
	)
	logger.Debug(row)
	logger.Debug(row.Err())
	var tmpLogin string
	err := row.Scan(&tmpLogin)
	logger.Debug(err)
	logger.Debug(tmpLogin)
	return err != sql.ErrNoRows
}
