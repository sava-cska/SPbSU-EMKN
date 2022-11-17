package accounts

type RevalidateChangePasswordCredentialsWrapper struct {
	RandomToken string `json:"random_token"`
	ExpiresIn   string `json:"expires_in"`
}

type RevalidateChangePasswordCredentialsResponse struct {
	Errors   *ErrorsUnion                                `json:"errors,omitempty"`
	Response *RevalidateChangePasswordCredentialsWrapper `json:"response,omitempty"`
}

func (response RevalidateChangePasswordCredentialsResponse) Bind() {}
