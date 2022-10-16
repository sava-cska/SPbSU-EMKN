package accounts

type RegisterWrapper struct {
	RandomToken string `json:"random_token,omitempty"`
	ExpiresIn   string `json:"expires_in,omitempty"`
}

type RegisterResponse struct {
	Errors   *ErrorsUnion     `json:"errors,omitempty"`
	Response *RegisterWrapper `json:"response,omitempty"`
}

func (response RegisterResponse) Bind() {}
