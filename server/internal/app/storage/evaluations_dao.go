package storage

import (
	"time"
)

type EvaluationsDAO struct {
	storage *Storage
}

type Record struct {
	UserId     string
	Evaluation string
	Result     *string
}

func (dao *EvaluationsDAO) Upsert(userId string, evaluation string, result *string) error {
	_, err := dao.storage.db.Exec(
		"INSERT INTO\n"+
			"evaluation_history (user_id, updated, evaluation, result)\n"+
			"VALUES ($1, $2, $3, $4)",
		userId,
		time.Now(),
		evaluation,
		result,
	)
	return err
}

func (dao *EvaluationsDAO) List(userId string) ([]Record, error) {
	rows, err := dao.storage.db.Query(
		"SELECT\n"+
			"	user_id,"+
			"	evaluation,"+
			"	result\n"+
			"FROM evaluation_history\n"+
			"WHERE user_id = $1\n"+
			"ORDER BY updated desc",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var eval Record
		if err := rows.Scan(&eval.UserId, &eval.Evaluation, &eval.Result); err != nil {
			return nil, err
		}
		records = append(records, eval)
	}
	return records, nil
}
