package accounts

type RevalidateChangePasswordCredentialsResponse struct {
	Errors *ErrorsUnion `json:"errors,omitempty"`
}

func (response RevalidateChangePasswordCredentialsResponse) Bind() {
}
