package accounts

type CommitChangePasswordRequest struct {
	ChangePasswordToken string `json:"change_password_token"`
	NewPassword         string `json:"new_password"`
}

func (request CommitChangePasswordRequest) Bind() {}
