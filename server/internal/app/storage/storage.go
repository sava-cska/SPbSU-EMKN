package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage struct {
	config         *Config
	db             *sql.DB
	evaluationsDao *EvaluationsDAO
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
	storage.db = db
	return nil
}

func (storage *Storage) Evaluations() *EvaluationsDAO {
	if storage.evaluationsDao == nil {
		storage.evaluationsDao = &EvaluationsDAO{
			storage: storage,
		}
	}
	return storage.evaluationsDao
}
