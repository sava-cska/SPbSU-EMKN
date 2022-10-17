package base

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/actions/accounts"
	"net/http"
	"strings"

	"github.com/gdexlab/go-render/render"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/error_handler"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/http_parser"
)

func HandleActionWithAuth[Req Request, Res Response](
	path string,
	businessLogicHandler func(*Req, *dependency.DependencyContext, ...any) (int, *Res),
	context *dependency.DependencyContext,
) {
	handleAction(path, businessLogicHandler, context, withAuth)
}

func HandleAction[Req Request, Res Response](
	path string,
	businessLogicHandler func(*Req, *dependency.DependencyContext, ...any) (int, *Res),
	context *dependency.DependencyContext,
) {
	handleAction(path, businessLogicHandler, context, EMPTY)
}

func handleAction[Req Request, Res Response](
	path string,
	businessLogicHandler func(*Req, *dependency.DependencyContext, ...any) (int, *Res),
	context *dependency.DependencyContext,
	middleware func(responseWriter http.ResponseWriter, httpRequest *http.Request, context *dependency.DependencyContext) (any, error),
) {
	handlerFunc := func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		context.Logger.Debugf("%s - Called URI %s", path, httpRequest.RequestURI)

		val, err := middleware(responseWriter, httpRequest, context)
		if err != nil {
			return
		}

		var request Req
		if errJSON := http_parser.ParseBody(interface{}(&request), httpRequest); errJSON != nil {
			error_handler.HandleError(context.Logger, responseWriter, http.StatusBadRequest, "Can't parse httpRequest.", errJSON)
			return
		}

		context.Logger.Debugf("%s\n\trequest: %s", path, render.AsCode(request))

		code, response := businessLogicHandler(&request, context, val)

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

var EMPTY = func(http.ResponseWriter, *http.Request, *dependency.DependencyContext) (any, error) { return "", nil }

func withAuth(responseWriter http.ResponseWriter,
	httpRequest *http.Request,
	context *dependency.DependencyContext) (any, error) {
	header := httpRequest.Header.Get("Authorization")
	if header == "" {
		err := errors.New("missing authorization header")
		error_handler.HandleError(context.Logger, responseWriter, http.StatusUnauthorized, err.Error(), err)
		return "", err
	}

	if !strings.HasPrefix(header, "Basic") {
		err := errors.New("unsupported authorization type")
		error_handler.HandleError(context.Logger, responseWriter, http.StatusUnauthorized, err.Error(), err)
		return "", err
	}

	authHeader := strings.TrimPrefix(header, "Basic ")
	context.Logger.Debugf(authHeader)
	parsed, err := base64.StdEncoding.DecodeString(authHeader)

	if err != nil {
		err := errors.New("wrong authorization format")
		error_handler.HandleError(context.Logger, responseWriter, http.StatusUnauthorized, err.Error(), err)
		return "", err
	}

	creds := strings.Split(string(parsed), ":")
	if len(creds) != 2 {
		err := errors.New("wrong authorization format")
		error_handler.HandleError(context.Logger, responseWriter, http.StatusUnauthorized, err.Error(), err)
		return "", err
	}
	login := creds[0]
	passwd := creds[1]

	isValid, err := accounts.ValidateUserCredentials(login, passwd, context.Storage)
	if err != nil {
		error_handler.HandleError(context.Logger, responseWriter, http.StatusInternalServerError, "Failed to validate credentials", err)
		return "", err
	}

	if !isValid {
		err := errors.New("wrong login or password")
		error_handler.HandleError(context.Logger, responseWriter, http.StatusUnauthorized, err.Error(), err)
		return "", err
	}
	return login, nil
}
