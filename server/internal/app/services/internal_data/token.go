package internal_data

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

const ResentCodeIn time.Duration = 60 * time.Second
const TokenTTL time.Duration = 30 * time.Minute

func GenerateToken() string {
	const tokenLength int = 20
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
