package storage

import (
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
)

type HomeworkDAO interface {
	GetAllHomeworks(courseId uint32) ([]models.HomeworkInDB, error)
}

type homeworkDAO struct {
	Storage *DbStorage
}

func (dao *homeworkDAO) GetAllHomeworks(courseId uint32) ([]models.HomeworkInDB, error) {
	res, err := dao.Storage.Db.Query(
		`SELECT id, name, deadline, score FROM homework_base
		 WHERE course_id = $1`, courseId,
	)
	if err != nil {
		return nil, err
	}

	homeworks := []models.HomeworkInDB{}
	for res.Next() {
		homework := models.HomeworkInDB{}
		var deadline_timestamp time.Time
		err = res.Scan(&homework.Id, &homework.Name, &deadline_timestamp, &homework.TotalScore)
		if err != nil {
			return nil, err
		}
		homework.CourseId = courseId
		homework.Deadline = deadline_timestamp.UnixMilli()
		homeworks = append(homeworks, homework)
	}

	return homeworks, nil
}
