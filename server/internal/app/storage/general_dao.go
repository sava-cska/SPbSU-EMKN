package storage

import "github.com/sava-cska/SPbSU-EMKN/internal/app/models"

type GeneralDAO interface {
	GetInfo() (*models.GeneralInfo, error)
}

type generalDAO struct {
	Storage *DbStorage
}

func (dao *generalDAO) GetInfo() (*models.GeneralInfo, error) {
	res := dao.Storage.Db.QueryRow(
		`SELECT current_period_id
               FROM general_base
               WHERE id = 0`,
	)

	info := models.GeneralInfo{}
	err := res.Scan(&info.CurrentPeriodId)
	return &info, err
}
