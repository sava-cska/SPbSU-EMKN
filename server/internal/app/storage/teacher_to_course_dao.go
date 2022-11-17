package storage

type TeacherToCourseDAO interface {
	GetTeachersByCourse(courseId uint32) ([]uint32, error)
}

type teacherToCourseDAO struct {
	Storage *DbStorage
}

func (dao *teacherToCourseDAO) GetTeachersByCourse(courseId uint32) ([]uint32, error) {
	rows, err := dao.Storage.Db.Query(
		`
		SELECT profile_id
		FROM teacher_to_course_base
		WHERE course_id = $1
		`,
		courseId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []uint32
	for rows.Next() {
		var teacher_profile uint32
		if errScan := rows.Scan(&teacher_profile); errScan != nil {
			return nil, errScan
		}
		teachers = append(teachers, teacher_profile)
	}
	if errRows := rows.Err(); errRows != nil {
		return nil, errRows
	}

	return teachers, nil
}
