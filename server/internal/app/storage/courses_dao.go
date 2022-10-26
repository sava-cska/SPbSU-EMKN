package storage

import (
	"database/sql"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"time"
)

type CoursesDAO struct {
	Storage *Storage
}

func (dao *CoursesDAO) GetPeriods() ([]*models.Period, error) {
	res, err := dao.Storage.Db.Query(
		`SELECT id, name
               FROM periods_base`,
	)
	if err != nil {
		return nil, err
	}

	periods := make([]*models.Period, 0)
	for res.Next() {
		period := models.Period{}
		err = res.Scan(&period.Id, &period.Text)
		if err != nil {
			return nil, err
		}
		periods = append(periods, &period)
	}

	return periods, nil
}

func (dao *CoursesDAO) GetDescription(courseId uint) (*string, error) {
	res := dao.Storage.Db.QueryRow(
		`SELECT short_description
               FROM courses_base
               WHERE id = $1`,
		courseId,
	)
	var description string
	err := res.Scan(&description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &description, nil
}

func (dao *CoursesDAO) GetDescriptionTimestamp(courseId uint) (*time.Time, error) {
	res := dao.Storage.Db.QueryRow(
		`SELECT description_timestamp
               FROM courses_base
               WHERE id = $1`,
		courseId,
	)
	var timestamp time.Time
	err := res.Scan(&timestamp)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &timestamp, nil
}
