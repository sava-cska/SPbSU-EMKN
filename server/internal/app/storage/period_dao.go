package storage

import "database/sql"

type PeriodDAO struct {
	Storage *Storage
}

func (dao *PeriodDAO) ExistPeriod(periodId uint32) (bool, error) {
	row := dao.Storage.Db.QueryRow(
		`
		SELECT id
		FROM period_base
		WHERE id = $1
		`,
		periodId,
	)
	var id uint32
	err := row.Scan(&id)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}
