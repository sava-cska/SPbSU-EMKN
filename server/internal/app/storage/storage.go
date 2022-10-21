package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage struct {
	config            *Config
	Db                *sql.DB
	userDAO           *UserDAO
	registrationDAO   *RegistrationDAO
	changePasswordDAO *ChangePasswordDAO
	userAvatarDAO     *UserAvatarDAO
	coursesDAO        *CoursesDAO
}

func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
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

func (storage *Storage) UserDAO() *UserDAO {
	if storage.userDAO == nil {
		storage.userDAO = &UserDAO{
			Storage: storage,
		}
	}
	return storage.userDAO
}

func (storage *Storage) RegistrationDAO() *RegistrationDAO {
	if storage.registrationDAO == nil {
		storage.registrationDAO = &RegistrationDAO{
			Storage: storage,
		}
	}
	return storage.registrationDAO
}

func (storage *Storage) ChangePasswordDAO() *ChangePasswordDAO {
	if storage.changePasswordDAO == nil {
		storage.changePasswordDAO = &ChangePasswordDAO{
			Storage: storage,
		}
	}
	return storage.changePasswordDAO
}

func (storage *Storage) UserAvatarDAO() *UserAvatarDAO {
	if storage.userAvatarDAO == nil {
		storage.userAvatarDAO = &UserAvatarDAO{
			Storage: storage,
		}
	}
	return storage.userAvatarDAO
}

func (storage *Storage) CoursesDAO() *CoursesDAO {
	if storage.coursesDAO == nil {
		storage.coursesDAO = &CoursesDAO{
			Storage: storage,
		}
	}
	return storage.coursesDAO
}
