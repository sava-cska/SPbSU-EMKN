package dependency

import (
	"github.com/gorilla/mux"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage/test_storage"
	"github.com/sirupsen/logrus"
	"io"
)

func NewTestContext() (*DependencyContext, *test_storage.TestDAO, *TestMailer) {
	db, dao := test_storage.New()
	mailer := TestMailer{Msgs: make(map[string][]*notifier.Message), MsgsSentCh: make(chan bool)}

	log := logrus.New()
	log.SetOutput(io.Discard)
	context := DependencyContext{
		Logger:  log,
		Router:  mux.NewRouter(),
		Storage: db,
		Mailer:  &mailer,
	}
	return &context, dao, &mailer
}

type TestMailer struct {
	Msgs       map[string][]*notifier.Message
	MsgsSentCh chan bool
}

func (m *TestMailer) SendEmail(receivers []string, message notifier.Message) error {
	for _, rec := range receivers {
		if _, ok := m.Msgs[rec]; !ok {
			m.Msgs[rec] = make([]*notifier.Message, 0)
		}
		m.Msgs[rec] = append(m.Msgs[rec], &message)
	}
	m.MsgsSentCh <- true
	return nil
}

func (m *TestMailer) WaitSingleMessageSent() {
	<-m.MsgsSentCh
}

func (m *TestMailer) GetLastEmail(email string) *notifier.Message {
	return m.Msgs[email][len(m.Msgs[email])-1]
}
