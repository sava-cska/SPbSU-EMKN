package accounts

type BeginChangePasswordRequest struct {
	Email string `json:"email"`
}

func (request BeginChangePasswordRequest) Bind() {}
