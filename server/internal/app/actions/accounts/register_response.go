package accounts

type Error struct {
	Code int `json:"code"`
}

type RegisterErrors struct {
	IllegalPassword     *Error `json:"illegal_password,omitempty"`
	IllegalLogin        *Error `json:"illegal_login,omitempty"`
	LoginIsNotAvailable *Error `json:"login_is_not_available,omitempty"`
}

type WrapResponse struct {
	RandomToken string `json:"random_token"`
	ExpiresIn   string `json:"expires_in"`
}

type RegisterResponse struct {
	Errors   RegisterErrors `json:"errors,omitempty"`
	Response WrapResponse   `json:"response"`
}
