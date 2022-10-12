package accounts

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (response LoginRequest) Bind() {
}
