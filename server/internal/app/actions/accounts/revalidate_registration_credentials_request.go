package accounts

type RevalidateRegistrationCredentialsRequest struct {
	Token string `json:"random_token"`
}
