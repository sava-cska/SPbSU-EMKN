package storage

import (
	"database/sql"
)

type Storage struct {
	config             *Config
	Db                 *sql.DB
	generalDAO         *GeneralDAO
	userDAO            *UserDAO
	registrationDAO    *RegistrationDAO
	changePasswordDAO  *ChangePasswordDAO
	userAvatarDAO      *UserAvatarDAO
	courseDAO          *CourseDAO
	periodDAO          *PeriodDAO
	teacherToCourseDAO *TeacherToCourseDAO
	studentToCourseDAO *StudentToCourseDAO
}

func New(config *Config) *Storage {
	s := &Storage{
		config: config,
	}
	s.generalDAO = &GeneralDAO{s}
	s.userDAO = &UserDAO{s}
	s.registrationDAO = &RegistrationDAO{s}
	s.changePasswordDAO = &ChangePasswordDAO{s}
	s.userAvatarDAO = &UserAvatarDAO{s}
	s.courseDAO = &CourseDAO{s}
	s.periodDAO = &PeriodDAO{s}
	s.teacherToCourseDAO = &TeacherToCourseDAO{s}
	s.studentToCourseDAO = &StudentToCourseDAO{s}
	return s
}

func (storage *Storage) Open() error {
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

func (storage *Storage) GeneralDAO() *GeneralDAO {
	return storage.generalDAO
}

func (storage *Storage) UserDAO() *UserDAO {
	return storage.userDAO
}

func (storage *Storage) RegistrationDAO() *RegistrationDAO {
	return storage.registrationDAO
}

func (storage *Storage) ChangePasswordDAO() *ChangePasswordDAO {
	return storage.changePasswordDAO
}

func (storage *Storage) UserAvatarDAO() *UserAvatarDAO {
	return storage.userAvatarDAO
}

func (storage *Storage) CourseDAO() *CourseDAO {
	return storage.courseDAO
}

func (storage *Storage) PeriodDAO() *PeriodDAO {
	return storage.periodDAO
}

func (storage *Storage) TeacherToCourseDAO() *TeacherToCourseDAO {
	return storage.teacherToCourseDAO
}

func (storage *Storage) StudentToCourseDAO() *StudentToCourseDAO {
	return storage.studentToCourseDAO
}
