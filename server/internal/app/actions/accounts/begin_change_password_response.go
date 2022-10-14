package accounts

type BeginChangePasswordWrapper struct {
	Token       string `json:"random_token,omitempty"`
	TimeExpired string `json:"expires_in,omitempty"`
}

type BeginChangePasswordResponse struct {
	Errors   *ErrorsUnion                `json:"errors,omitempty"`
	Response *BeginChangePasswordWrapper `json:"response,omitempty"`
}

func (response BeginChangePasswordResponse) Bind() {}
