package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {
	GeneralDAO() GeneralDAO
	UserDAO() UserDAO
	RegistrationDAO() RegistrationDAO
	ChangePasswordDAO() ChangePasswordDAO
	UserAvatarDAO() UserAvatarDAO
	CourseDAO() CourseDAO
	PeriodDAO() PeriodDAO
	TeacherToCourseDAO() TeacherToCourseDAO
	StudentToCourseDAO() StudentToCourseDAO
}

type DbStorage struct {
	config             *Config
	Db                 *sql.DB
	generalDAO         GeneralDAO
	userDAO            UserDAO
	registrationDAO    RegistrationDAO
	changePasswordDAO  ChangePasswordDAO
	userAvatarDAO      UserAvatarDAO
	courseDAO          CourseDAO
	periodDAO          PeriodDAO
	teacherToCourseDAO TeacherToCourseDAO
	studentToCourseDAO StudentToCourseDAO
}

func New(config *Config) *DbStorage {
	s := &DbStorage{
		config: config,
	}
	s.generalDAO = &generalDAO{s}
	s.userDAO = &userDAO{s}
	s.registrationDAO = &registrationDAO{s}
	s.changePasswordDAO = &changePasswordDAO{s}
	s.userAvatarDAO = &userAvatarDAO{s}
	s.courseDAO = &courseDAO{s}
	s.periodDAO = &periodDAO{s}
	s.teacherToCourseDAO = &teacherToCourseDAO{s}
	s.studentToCourseDAO = &studentToCourseDAO{s}
	return s
}

func (storage *DbStorage) Open() error {
	db, err := sql.Open("postgres", storage.config.DatabaseURL)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	storage.Db = db
	return nil
}

func (storage *DbStorage) GeneralDAO() GeneralDAO {
	return storage.generalDAO
}

func (storage *DbStorage) UserDAO() UserDAO {
	return storage.userDAO
}

func (storage *DbStorage) RegistrationDAO() RegistrationDAO {
	return storage.registrationDAO
}

func (storage *DbStorage) ChangePasswordDAO() ChangePasswordDAO {
	return storage.changePasswordDAO
}

func (storage *DbStorage) UserAvatarDAO() UserAvatarDAO {
	return storage.userAvatarDAO
}

func (storage *DbStorage) CourseDAO() CourseDAO {
	return storage.courseDAO
}

func (storage *DbStorage) PeriodDAO() PeriodDAO {
	return storage.periodDAO
}

func (storage *DbStorage) TeacherToCourseDAO() TeacherToCourseDAO {
	return storage.teacherToCourseDAO
}

func (storage *DbStorage) StudentToCourseDAO() StudentToCourseDAO {
	return storage.studentToCourseDAO
}
