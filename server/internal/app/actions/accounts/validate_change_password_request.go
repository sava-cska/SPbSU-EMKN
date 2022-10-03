package accounts

type ValidateChangePasswordRequest struct {
	RandomToken      string `json:"random_token"`
	VerificationCode string `json:"verification_code"`
}
