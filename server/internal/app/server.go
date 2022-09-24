package server

import (
	"github.com/sava-cska/SPbSU-Calculator/internal/app/actions/evaluations"
	"github.com/sava-cska/SPbSU-Calculator/internal/app/storage"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config  *Config
	logger  *logrus.Logger
	router  *mux.Router
	storage *storage.Storage
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (server *Server) Start() error {
	if err := server.configureLogger(); err != nil {
		return err
	}
	if err := server.configureStorage(); err != nil {
		return err
	}
	server.configureRouter()
	server.logger.Info("Server is up")
	return http.ListenAndServe(server.config.BindAddress, server.router)
}

func (server *Server) configureLogger() error {
	level, err := logrus.ParseLevel(server.config.LogLevel)
	if err != nil {
		return err
	}
	server.logger.SetLevel(level)
	server.logger.SetFormatter(configureLogFormatter())
	return nil
}

func configureLogFormatter() *logrus.TextFormatter {
	formatter := new(logrus.TextFormatter)
	formatter.FullTimestamp = true
	return formatter
}

func (server *Server) configureStorage() error {
	localStorage := storage.New(server.config.Storage)
	if err := localStorage.Open(); err != nil {
		return err
	}
	server.storage = localStorage
	return nil
}

func (server *Server) configureRouter() {
	server.router.HandleFunc("/evaluations/list", evaluations.HandleEvaluationsList(server.logger, server.storage))
	server.router.HandleFunc("/evaluations/calculate", evaluations.HandleEvaluationsCalculate(server.logger, server.storage))
}
