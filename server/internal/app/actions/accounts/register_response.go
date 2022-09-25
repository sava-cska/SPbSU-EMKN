package accounts

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/actions"
)

type RegisterResponse struct {
	actions.EmknResponse
	RandomToken string `json:"random_token"`
	ExpiresIn   string `json:"expires_in"`
}
