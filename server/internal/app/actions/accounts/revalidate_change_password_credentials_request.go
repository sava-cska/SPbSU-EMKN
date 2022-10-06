package accounts

type RevalidateChangePasswordCredentialsRequest struct {
	RandomToken string `json:"random_token"`
}
