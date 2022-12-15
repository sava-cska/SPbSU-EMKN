package core

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	_ "github.com/sava-cska/SPbSU-EMKN/docs"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/actions/accounts"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/actions/base"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/actions/courses"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/actions/profiles"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/event_queue"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
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
	rabbitLogin := os.Getenv("RABBIT_LOGIN")
	rabbitPassword := os.Getenv("RABBIT_PASSWORD")

	if err := server.configureLogger(); err != nil {
		return err
	}
	if err := server.configureStorage(); err != nil {
		return err
	}
	if err := server.configureQueue(rabbitLogin, rabbitPassword); err != nil {
		return err
	}

	// important to have configured other entities before configure router
	server.configureRouter()
	server.context.Router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
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

func (server *Server) configureQueue(rabbitLogin, rabbitPassword string) error {
	queue, err := event_queue.New(server.config.EventQueueConfig, rabbitLogin, rabbitPassword)
	server.context.EventQueue = queue
	return err
}

func (server *Server) configureRouter() {
	base.HandleAction("/accounts/register", accounts.HandleAccountsRegister, server.context)
	base.HandleAction("/accounts/validate_email", accounts.HandleAccountsValidateEmail, server.context)
	base.HandleAction("/accounts/login", accounts.HandleAccountsLogin, server.context)
	base.HandleAction("/accounts/begin_change_password", accounts.HandleAccountsBeginChangePassword, server.context)
	base.HandleAction("/accounts/validate_change_password", accounts.HandleAccountsValidateChangePassword, server.context)
	base.HandleAction("/accounts/commit_change_password", accounts.HandleAccountsCommitChangePassword, server.context)
	base.HandleAction("/accounts/revalidate_registration_credentials", accounts.HandleAccountsRevalidateRegistrationCredentials,
		server.context)
	base.HandleAction("/accounts/revalidate_change_password_credentials", accounts.HandleAccountsRevalidateChangePasswordCredentials,
		server.context)
	base.HandleActionWithAuth("/profiles/get", profiles.HandleProfilesGet, server.context)
	base.HandleActionWithAuth("/profiles/load_image", profiles.HandleProfilesLoadImage, server.context)
	base.HandleActionWithAuth("/courses/periods", courses.HandleCoursesPeriods, server.context)
	base.HandleActionWithAuth("/courses/description", courses.HandleCoursesDescription, server.context)
	base.HandleActionWithAuth("/courses/description_ping", courses.HandleCoursesDescriptionPing, server.context)
	base.HandleActionWithAuth("/courses/list", courses.HandleCoursesList, server.context)
	base.HandleActionWithAuth("/courses/enroll", courses.HandleCoursesEnroll, server.context)
	base.HandleActionWithAuth("/courses/unenroll", courses.HandleCoursesUnEnroll, server.context)
	base.HandleActionWithAuth("/courses/get_homeworks", courses.HandleCoursesGetHomeworks, server.context)
}
