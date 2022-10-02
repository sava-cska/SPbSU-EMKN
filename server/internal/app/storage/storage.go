package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage struct {
	config          *Config
	Db              *sql.DB
	userDao         *UserDAO
	registrationDao *RegistrationDAO
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
	if storage.userDao == nil {
		storage.userDao = &UserDAO{
			Storage: storage,
		}
	}
	return storage.userDao
}

func (storage *Storage) RegistrationDAO() *RegistrationDAO {
	if storage.registrationDao == nil {
		storage.registrationDao = &RegistrationDAO{
			Storage: storage,
		}
	}
	return storage.registrationDao
}
