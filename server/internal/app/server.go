package server

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/utils"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/actions/accounts"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config  *Config
	logger  *logrus.Logger
	router  *mux.Router
	storage *storage.Storage
	mailer  *notifier.Mailer
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (server *Server) Start() error {
	rand.Seed(time.Now().UnixNano())
	EmknCourseMail := os.Getenv("EMKN_COURSE_MAIL")
	EmknCoursePassword := os.Getenv("EMKN_COURSE_PASSWORD")

	if err := server.configureLogger(); err != nil {
		return err
	}
	if err := server.configureStorage(); err != nil {
		return err
	}
	server.configureMailing(EmknCourseMail, EmknCoursePassword)

	// important to have configured other entities before configure router
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
	localStorage := storage.New(server.config.StorageConfig)
	if err := localStorage.Open(); err != nil {
		return err
	}
	server.storage = localStorage
	return nil
}

func (server *Server) configureMailing(EmknCourseMail, EmknCoursePassword string) {
	mailer := notifier.New(server.config.NotifierConfig, EmknCourseMail, EmknCoursePassword)
	server.mailer = mailer
}

func (server *Server) configureRouter() {
	server.router.HandleFunc("/accounts/register",
		accounts.HandleAccountsRegister(server.logger, server.storage, server.mailer))
	server.router.HandleFunc("/accounts/validate_email",
		accounts.HandleAccountsValidateEmail(server.logger, server.storage))
	server.router.HandleFunc("/accounts/login",
		accounts.HandleAccountsLogin(server.logger, server.storage))
	server.router.HandleFunc("/accounts/begin_change_password",
		accounts.HandleAccountsChangePwd(server.logger, server.storage, server.mailer))
	server.router.HandleFunc("/accounts/validate_change_password",
		accounts.HandleValidateChangePassword(server.logger, server.storage))
	server.router.HandleFunc("/accounts/commit_change_password",
		accounts.HandleAccountsCommitPwdChange(server.logger, server.storage))
	server.router.HandleFunc("/accounts/revalidate_registration_credentials",
		accounts.HandleAccountsRevalidateRegistrationCredentials(server.logger, server.storage, server.mailer))
	server.router.HandleFunc("/accounts/revalidate_change_password_credentials",
		accounts.HandleRevalidateChangePasswordCredentials(server.logger, server.storage, server.mailer))
}

// used before all handlers that require user authorization
func (server *Server) withAuth(handlerFunc http.HandlerFunc, logger *logrus.Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		header := request.Header.Get("Authorization")
		if header == "" {
			utils.HandleError(logger, writer, http.StatusUnauthorized, "Missing authorization header", nil)
			return
		}

		if !strings.HasPrefix(header, "Basic") {
			utils.HandleError(logger, writer, http.StatusUnauthorized, "Unsupported authorization type", nil)
			return
		}

		authHeader := strings.TrimPrefix(header, "Basic ")
		creds := strings.Split(authHeader, ":")
		if len(creds) != 2 {
			utils.HandleError(logger, writer, http.StatusUnauthorized, "Wrong authorization format", nil)
			return
		}
		login := creds[0]
		passwd := creds[1]

		isValid, err := accounts.ValidateUserCredentials(login, passwd, server.storage)
		if err != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Failed to validate credentials", err)
			return
		}

		if !isValid {
			utils.HandleError(logger, writer, http.StatusUnauthorized, "Wrong login or password", nil)
			return
		}

		handlerFunc(writer, request)
	}
}
