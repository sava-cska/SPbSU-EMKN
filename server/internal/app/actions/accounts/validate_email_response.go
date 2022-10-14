package accounts

type ValidateEmailResponse struct {
	Errors *ErrorsUnion `json:"errors,omitempty"`
}

func (response ValidateEmailResponse) Bind() {}
