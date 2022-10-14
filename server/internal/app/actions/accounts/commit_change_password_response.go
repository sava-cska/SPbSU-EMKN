package accounts

type CommitChangePasswordResponse struct {
	Errors *ErrorsUnion `json:"errors,omitempty"`
}

func (response CommitChangePasswordResponse) Bind() {}
