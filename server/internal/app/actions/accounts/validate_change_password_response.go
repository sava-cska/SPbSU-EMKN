package accounts

type ValidateChangePasswordResponseWrapper struct {
	ChangePasswordToken string `json:"change_password_token,omitempty"`
}

type ValidateChangePasswordResponse struct {
	Response *ValidateChangePasswordResponseWrapper `json:"response,omitempty"`
	Errors   *ErrorsUnion                           `json:"errors,omitempty"`
}

func (response ValidateChangePasswordResponse) Bind() {}
