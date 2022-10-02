package accounts

type Error struct {
	Code int `json:"code"`
}

type RegisterErrors struct {
	IllegalPassword     *Error `json:"illegal_password,omitempty"`
	IllegalLogin        *Error `json:"illegal_login,omitempty"`
	IllegalEmail        *Error `json:"illegal_email,omitempty"`
	LoginIsNotAvailable *Error `json:"login_is_not_available,omitempty"`
}

type WrapResponse struct {
	RandomToken string `json:"random_token,omitempty"`
	ExpiresIn   string `json:"expires_in,omitempty"`
}

type RegisterResponse struct {
	Errors   *RegisterErrors `json:"errors,omitempty"`
	Response *WrapResponse   `json:"response,omitempty"`
}
