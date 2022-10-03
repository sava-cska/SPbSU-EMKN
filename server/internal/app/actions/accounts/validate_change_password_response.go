package accounts

type ValidateChangePasswordResponse struct {
	ChangePasswordToken string       `json:"change_password_token,omitempty"`
	Errors              *ErrorsUnion `json:"errors,omitempty"`
}
