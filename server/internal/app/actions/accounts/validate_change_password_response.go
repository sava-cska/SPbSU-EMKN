package accounts

type ValidateChangePasswordResponse struct {
	ChangePasswordToken string       `json:"change_password_token"`
	Errors              *ErrorsUnion `json:"errors,omitempty"`
}

func (response ValidateChangePasswordResponse) Bind() {}
