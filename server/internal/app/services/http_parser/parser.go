package http_parser

import (
	"encoding/json"
	"io"
	"net/http"
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
