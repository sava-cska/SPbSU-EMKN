package accounts

type RevalidateRegistrationCredentialsWrapper struct {
	RandomToken string `json:"random_token"`
	ExpiresIn   string `json:"expires_in"`
}

type RevalidateRegistrationCredentialsResponse struct {
	Errors   *ErrorsUnion                              `json:"errors,omitempty"`
	Response *RevalidateRegistrationCredentialsWrapper `json:"response,omitempty"`
}

func (request RevalidateRegistrationCredentialsResponse) Bind() {}
