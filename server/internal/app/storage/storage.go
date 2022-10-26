package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage struct {
	config            *Config
	Db                *sql.DB
	generalDAO        *GeneralDAO
	userDAO           *UserDAO
	registrationDAO   *RegistrationDAO
	changePasswordDAO *ChangePasswordDAO
	userAvatarDAO     *UserAvatarDAO
	coursesDAO        *CoursesDAO
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
	s.coursesDAO = &CoursesDAO{s}

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

func (storage *Storage) CoursesDAO() *CoursesDAO {
	return storage.coursesDAO
}
