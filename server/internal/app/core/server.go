package core

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/actions/accounts"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/actions/base"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config  *Config
	context *dependency.DependencyContext
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		context: &dependency.DependencyContext{
			Logger: logrus.New(),
			Router: mux.NewRouter(),
		},
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
	server.context.Logger.Info("Server is up")
	return http.ListenAndServe(server.config.BindAddress, server.context.Router)
}

func (server *Server) configureLogger() error {
	level, err := logrus.ParseLevel(server.config.LogLevel)
	if err != nil {
		return err
	}
	server.context.Logger.SetLevel(level)
	server.context.Logger.SetFormatter(configureLogFormatter())
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
	server.context.Storage = localStorage
	return nil
}

func (server *Server) configureMailing(EmknCourseMail, EmknCoursePassword string) {
	mailer := notifier.New(server.config.NotifierConfig, EmknCourseMail, EmknCoursePassword)
	server.context.Mailer = mailer
}

func (server *Server) configureRouter() {
	base.HandleAction("/accounts/register", accounts.HandleAccountsRegister, server.context)

	//server.router.HandleFunc("/accounts/validate_email",
	//	accounts.HandleAccountsValidateEmail(server.logger, server.storage))

	base.HandleAction("/accounts/login", accounts.HandleAccountsLogin, server.context)

	//server.router.HandleFunc("/accounts/begin_change_password",
	//	accounts.HandleAccountsChangePwd(server.logger, server.storage, server.mailer))

	base.HandleAction("/accounts/validate_change_password", accounts.HandleValidateChangePassword, server.context)

	//server.router.HandleFunc("/accounts/commit_change_password",
	//	accounts.HandleAccountsCommitPwdChange(server.logger, server.storage))
	//server.router.HandleFunc("/accounts/revalidate_registration_credentials",
	//	accounts.HandleAccountsRevalidateRegistrationCredentials(server.logger, server.storage, server.mailer))

	base.HandleAction("/accounts/revalidate_change_password_credentials", accounts.HandleRevalidateChangePasswordCredentials, server.context)
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

		isValid, err := accounts.ValidateUserCredentials(login, passwd, server.context.Storage)
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
