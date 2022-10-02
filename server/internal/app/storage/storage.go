package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage struct {
	config          *Config
	Db              *sql.DB
	userDao         *UserDAO
	registrationDao *RegistrationDao
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
func (storage *Storage) UserDao() *UserDAO {
	if storage.userDao == nil {
		storage.userDao = &UserDAO{
			Storage: storage,
		}
	}
	return storage.userDao
}

func (storage *Storage) RegistrationDao() *RegistrationDao {
	if storage.registrationDao == nil {
		storage.registrationDao = &RegistrationDao{
			Storage: storage,
		}
	}
	return storage.registrationDao
}
