package accounts

type LoginResponse struct {
	Errors *ErrorsUnion `json:"errors,omitempty"`
}
