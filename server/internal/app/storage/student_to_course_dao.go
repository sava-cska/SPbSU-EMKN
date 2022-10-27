package storage

import "database/sql"

type StudentToCourseDAO struct {
	Storage *Storage
}

func (dao *StudentToCourseDAO) ExistRecord(profileId uint32, courseId uint32) (bool, error) {
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
