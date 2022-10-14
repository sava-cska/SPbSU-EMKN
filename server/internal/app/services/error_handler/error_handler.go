package error_handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func HandleError(logger *logrus.Logger, writer http.ResponseWriter, code int, reason string, err error) {
	logger.Error(reason, err)
	writer.WriteHeader(code)
	writer.Write([]byte(reason))
}
