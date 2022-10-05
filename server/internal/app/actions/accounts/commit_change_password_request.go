package accounts

type CommitPwdChangeRequest struct {
	ChangePwdToken string `json:"change_password_token"`
	NewPassword    string `json:"new_password"`
}
