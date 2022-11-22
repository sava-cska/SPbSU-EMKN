package core

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage/test_storage"
	"net/http/httptest"
	"strings"
)

func NewTest() (*TestServer, *test_storage.TestDAO, *dependency.TestMailer) {
	ctx, db, mailer := dependency.NewTestContext()
	server := &Server{
		context: ctx,
	}
	server.configureRouter()

	return &TestServer{router: ctx.Router}, db, mailer
}

type TestServer struct {
	router *mux.Router
}

func (t *TestServer) SendRequest(endpoint string, body string) (int, string) {
	req := httptest.NewRequest("POST", fmt.Sprintf("http://www.emkn.ru%s", endpoint), strings.NewReader(body))
	w := httptest.NewRecorder()
	t.router.ServeHTTP(w, req)

	return w.Code, w.Body.String()
}
