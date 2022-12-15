package storage

import "database/sql"

type PassedHomeworkDAO interface {
	CheckUserPassHomework(userId uint32, homeworkId uint32) (bool, error)
}

type passedHomeworkDAO struct {
	Storage *DbStorage
}

func (dao *passedHomeworkDAO) CheckUserPassHomework(userId uint32, homeworkId uint32) (bool, error) {
	res := dao.Storage.Db.QueryRow(
		`SELECT user_id
         FROM passed_homework_base
         WHERE user_id = $1 AND homework_id = $2`,
		userId, homeworkId,
	)
	var id uint32
	err := res.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
