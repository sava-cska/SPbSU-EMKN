package storage

import (
	"database/sql"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
)

type CourseDAO struct {
	Storage *Storage
}

func (dao *CourseDAO) GetPeriods() ([]*models.Period, error) {
	res, err := dao.Storage.Db.Query(
		`SELECT id, name
               FROM period_base`,
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

func (dao *CourseDAO) GetDescription(courseId uint32) (*string, error) {
	res := dao.Storage.Db.QueryRow(
		`SELECT short_description
               FROM course_base
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

func (dao *CourseDAO) GetDescriptionTimestamp(courseId uint32) (*time.Time, error) {
	res := dao.Storage.Db.QueryRow(
		`SELECT description_timestamp
               FROM course_base
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

func (dao *CourseDAO) GetCoursesByPeriod(periodId uint32) ([]*models.CourseInDB, error) {
	rows, err := dao.Storage.Db.Query(
		`
		SELECT id, title, short_description
		FROM course_base
		WHERE period_id = $1
		`,
		periodId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.CourseInDB
	for rows.Next() {
		var course models.CourseInDB
		if errScan := rows.Scan(&course.Id, &course.Title, &course.ShortDescription); errScan != nil {
			return nil, errScan
		}
		courses = append(courses, &course)
	}
	if errRows := rows.Err(); errRows != nil {
		return nil, errRows
	}

	return courses, nil
}

func (dao *CourseDAO) ExistCourse(courseId uint32) (bool, error) {
	row := dao.Storage.Db.QueryRow(
		`
		SELECT id
		FROM course_base
		WHERE id = $1
		`,
		courseId,
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
