package dependency

import (
	"github.com/gorilla/mux"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage/test_storage"
	"github.com/sirupsen/logrus"
	"io"
)

func NewTestContext() (*DependencyContext, *test_storage.TestDAO, *TestMailer) {
	db, dao := test_storage.New()
	mailer := TestMailer{Msgs: make(map[string][]*models.Message), MsgsSentCh: make(chan bool)}

	log := logrus.New()
	log.SetOutput(io.Discard)
	context := DependencyContext{
		Logger:  log,
		Router:  mux.NewRouter(),
		Storage: db,
	}
	return &context, dao, &mailer
}

type TestMailer struct {
	Msgs       map[string][]*models.Message
	MsgsSentCh chan bool
}

func (m *TestMailer) SendEmail(message models.Message) error {
	for _, rec := range message.Receivers {
		if _, ok := m.Msgs[rec]; !ok {
			m.Msgs[rec] = make([]*models.Message, 0)
		}
		m.Msgs[rec] = append(m.Msgs[rec], &message)
	}
	m.MsgsSentCh <- true
	return nil
}

func (m *TestMailer) WaitSingleMessageSent() {
	<-m.MsgsSentCh
}

func (m *TestMailer) GetLastEmail(email string) *models.Message {
	return m.Msgs[email][len(m.Msgs[email])-1]
}
