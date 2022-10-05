package utils

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const ResentCodeIn time.Duration = 60 * time.Second
const TokenTTL time.Duration = 30 * time.Minute

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

func GenerateToken() string {
	const tokenLength int = 20
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func GenerateVerificationCode() string {
	const verificationCodeLength int = 6
	code := strings.Builder{}
	for i := 0; i < verificationCodeLength; i++ {
		code.WriteString(strconv.Itoa(rand.Intn(10)))
	}
	return code.String()
}
