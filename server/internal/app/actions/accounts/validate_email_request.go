package accounts

type ValidateEmailRequest struct {
	VerificationCode string `json:"verification_code"`
	Token            string `json:"random_token"`
}

func (request ValidateEmailRequest) Bind() {}
