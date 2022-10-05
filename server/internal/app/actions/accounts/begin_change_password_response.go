package accounts

type ChangePwdWrapper struct {
	Token       string `json:"random_token"`
	TimeExpired string `json:"expires_in"`
}

type ChangePwdResponse struct {
	Errors   *ErrorsUnion      `json:"errors,omitempty"`
	Response *ChangePwdWrapper `json:"response,omitempty"`
}
