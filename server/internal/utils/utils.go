package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func ParseBody(value any, request *http.Request) error {
	data, errReq := io.ReadAll(request.Body)
	if errReq != nil {
		return errReq
	}
	if errJSON := json.Unmarshal(data, value); errJSON != nil {
		return errJSON
	}

	return nil
}

func HandleError(logger *logrus.Logger, writer http.ResponseWriter, code int, reason string, err error) {
	logger.Error(reason, err)
	writer.WriteHeader(code)
	writer.Write([]byte(reason))
}
