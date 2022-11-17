package test_storage

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
)

type TestStorage struct {
	dao *TestDAO
}

func New() (*TestStorage, *TestDAO) {
	dao := TestDAO{
		LoginToUser:           make(map[string]*models.User),
		TokenToRegistration:   make(map[string]*TestRegistrationData),
		CurrentPeriodId:       0,
		Periods:               make(map[uint32]*models.Period),
		Courses:               make(map[uint32]*TestCourseData),
		TokenToChangePassword: make(map[string]*TestChangePasswordData),
		UserAvatars:           make(map[uint32]*models.Profile),
	}
	return &TestStorage{dao: &dao}, &dao
}

func (storage *TestStorage) GeneralDAO() storage.GeneralDAO {
	return storage.dao
}

func (storage *TestStorage) UserDAO() storage.UserDAO {
	return storage.dao
}

func (storage *TestStorage) RegistrationDAO() storage.RegistrationDAO {
	return storage.dao
}

func (storage *TestStorage) ChangePasswordDAO() storage.ChangePasswordDAO {
	return storage.dao
}

func (storage *TestStorage) UserAvatarDAO() storage.UserAvatarDAO {
	return storage.dao
}

func (storage *TestStorage) CourseDAO() storage.CourseDAO {
	return storage.dao
}

func (storage *TestStorage) PeriodDAO() storage.PeriodDAO {
	return storage.dao
}

func (storage *TestStorage) TeacherToCourseDAO() storage.TeacherToCourseDAO {
	return storage.dao
}

func (storage *TestStorage) StudentToCourseDAO() storage.StudentToCourseDAO {
	return storage.dao
}
