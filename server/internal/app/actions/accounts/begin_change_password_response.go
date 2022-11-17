package accounts

type BeginChangePasswordWrapper struct {
	Token       string `json:"random_token"`
	TimeExpired string `json:"expires_in"`
}

type BeginChangePasswordResponse struct {
	Errors   *ErrorsUnion                `json:"errors,omitempty"`
	Response *BeginChangePasswordWrapper `json:"response,omitempty"`
}

func (response BeginChangePasswordResponse) Bind() {}
