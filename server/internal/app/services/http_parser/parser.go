package http_parser

import (
	"encoding/json"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"io"
	"net/http"
)

func ParseBody(value any, request *http.Request, context *dependency.DependencyContext) error {
	data, errReq := io.ReadAll(request.Body)
	if errReq == io.EOF {
		return nil
	}
	if errReq != nil {
		context.Logger.Errorf("ParseBody readall error: %s", errReq)
		return errReq
	}
	if errJSON := json.Unmarshal(data, value); errJSON != nil {
		context.Logger.Errorf("ParseBody unmarshal error: %s", errJSON)
		return nil
	}
	return nil
}
