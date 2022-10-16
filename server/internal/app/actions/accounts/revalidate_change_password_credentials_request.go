package accounts

type RevalidateChangePasswordCredentialsRequest struct {
	RandomToken string `json:"random_token"`
}

func (response RevalidateChangePasswordCredentialsRequest) Bind() {}
