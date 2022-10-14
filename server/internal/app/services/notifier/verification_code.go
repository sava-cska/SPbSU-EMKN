package notifier

import (
	"math/rand"
	"strconv"
	"strings"
)

func GenerateVerificationCode() string {
	const verificationCodeLength int = 6
	code := strings.Builder{}
	for i := 0; i < verificationCodeLength; i++ {
		code.WriteString(strconv.Itoa(rand.Intn(10)))
	}
	return code.String()
}
