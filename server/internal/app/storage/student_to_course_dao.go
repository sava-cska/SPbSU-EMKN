package storage

import "database/sql"

type StudentToCourseDAO interface {
	ExistRecord(profileId uint32, courseId uint32) (bool, error)
	AddRecord(profileId uint32, courseId uint32) error
	DeleteRecord(profileId uint32, courseId uint32) error
}

type studentToCourseDAO struct {
	Storage *DbStorage
}

func (dao *studentToCourseDAO) ExistRecord(profileId uint32, courseId uint32) (bool, error) {
	row := dao.Storage.Db.QueryRow(
		`
		SELECT profile_id
		FROM student_to_course_base
		WHERE profile_id = $1 AND course_id = $2
		`,
		profileId,
		courseId,
	)
	var profile uint32
	err := row.Scan(&profile)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func (dao *studentToCourseDAO) AddRecord(profileId uint32, courseId uint32) error {
	_, err := dao.Storage.Db.Exec(
		`
		INSERT INTO student_to_course_base
		(profile_id, course_id)
		VALUES ($1, $2)
		`,
		profileId, courseId)
	return err
}

func (dao *studentToCourseDAO) DeleteRecord(profileId uint32, courseId uint32) error {
	_, err := dao.Storage.Db.Exec(
		`
		DELETE FROM student_to_course_base
		WHERE profile_id = $1 AND course_id = $2
		`,
		profileId,
		courseId,
	)
	return err
}
