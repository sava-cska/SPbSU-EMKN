package accounts

type LoginResponse struct {
	Errors *ErrorsUnion `json:"errors,omitempty"`
}

func (response LoginResponse) Bind() {
}
