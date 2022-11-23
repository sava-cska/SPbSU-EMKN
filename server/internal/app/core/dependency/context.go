package dependency

import (
	"github.com/gorilla/mux"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/event_queue"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sirupsen/logrus"
)

type DependencyContext struct {
	Logger     *logrus.Logger
	Router     *mux.Router
	Storage    storage.Storage
	EventQueue event_queue.EventQueue
}
