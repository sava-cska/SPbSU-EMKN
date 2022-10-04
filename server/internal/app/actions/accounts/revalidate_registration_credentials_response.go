package accounts

type RevalidateRegistrationCredentialsWrapper struct {
	RandomToken string `json:"random_token,omitempty"`
	ExpiresIn   string `json:"expires_in,omitempty"`
}

type RevalidateRegistrationCredentialsResponse struct {
	Errors   *ErrorsUnion                              `json:"errors,omitempty"`
	Response *RevalidateRegistrationCredentialsWrapper `json:"response,omitempty"`
}
