package base

import (
	"encoding/json"
	"net/http"

	"github.com/gdexlab/go-render/render"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/error_handler"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/http_parser"
)

func HandleAction[Req Request, Res Response](
	path string,
	businessLogicHandler func(*Req, *dependency.DependencyContext) (int, *Res),
	context *dependency.DependencyContext,
) {
	handlerFunc := func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		context.Logger.Debugf("%s - Called URI %s", path, httpRequest.RequestURI)

		var request Req
		if errJSON := http_parser.ParseBody(interface{}(&request), httpRequest); errJSON != nil {
			error_handler.HandleError(context.Logger, responseWriter, http.StatusBadRequest, "Can't parse httpRequest.", errJSON)
			return
		}

		context.Logger.Debugf("%s\n\trequest: %s", path, render.AsCode(request))

		code, response := businessLogicHandler(&request, context)

		respJSON, errRespJSON := json.Marshal(response)
		if errRespJSON != nil {
			error_handler.HandleError(context.Logger, responseWriter, http.StatusInternalServerError,
				"Can't create JSON object from data.", errRespJSON)
			return
		}
		context.Logger.Debugf("%s\n\tcode: %d\n\tresponse: %s", path, code, render.AsCode(response))

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(code)
		responseWriter.Write(respJSON)
	}
	context.Router.HandleFunc(path, handlerFunc)
}
