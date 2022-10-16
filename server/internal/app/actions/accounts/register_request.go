package accounts

type RegisterRequest struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (request RegisterRequest) Bind() {}
