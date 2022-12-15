package storage

import "database/sql"

type CheckedHomeworkDAO interface {
	ScoreForHomework(userId uint32, homeworkId uint32) (bool, int, error)
}

type checkedHomeworkDAO struct {
	Storage *DbStorage
}

func (dao *checkedHomeworkDAO) ScoreForHomework(userId uint32, homeworkId uint32) (bool, int, error) {
	res := dao.Storage.Db.QueryRow(
		`SELECT user_score
         FROM checked_homework_base
         WHERE user_id = $1 AND homework_id = $2`,
		userId, homeworkId,
	)
	var score int
	err := res.Scan(&score)
	if err == sql.ErrNoRows {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	return true, score, nil
}
