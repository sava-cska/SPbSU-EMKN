package evaluations

import (
	"encoding/json"
	"fmt"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Token struct {
	Body string
	Type string
}

type WrapperTokens struct {
	Tokens []Token
}

func toString(tokens []Token) string {
	var sb strings.Builder
	for _, token := range tokens {
		sb.WriteString(token.Body)
	}
	return sb.String()
}

func encode(tokens []Token, logger *logrus.Logger, writer http.ResponseWriter) string {
	encodeTokens := WrapperTokens{Tokens: tokens}
	jsonTokens, err := json.Marshal(encodeTokens)
	if err != nil {
		utils.HandleError(logger, writer, http.StatusInternalServerError,
			fmt.Sprintf("Can't create json from expression %s", toString(tokens)), err)
		return ""
	}
	return string(jsonTokens)
}

func decode(encodedToken string, logger *logrus.Logger, writer http.ResponseWriter) []Token {
	var decodeTokens WrapperTokens
	if err := json.Unmarshal([]byte(encodedToken), &decodeTokens); err != nil {
		utils.HandleError(logger, writer, http.StatusInternalServerError,
			fmt.Sprintf("Can't parse json from database %s", encodedToken), err)
		return []Token{}
	}
	return decodeTokens.Tokens
}
