package storage

import "github.com/sava-cska/SPbSU-EMKN/internal/app/models"

type GeneralDAO struct {
	Storage *Storage
}

func (dao *GeneralDAO) GetInfo() (*models.GeneralInfo, error) {
	res := dao.Storage.Db.QueryRow(
		`SELECT current_period_id
               FROM general_base
               WHERE id = 0`,
	)

	info := models.GeneralInfo{}
	err := res.Scan(&info.CurrentPeriodId)
	return &info, err
}