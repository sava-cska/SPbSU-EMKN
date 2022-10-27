package accounts

type LoginResponse struct {
	Errors    *ErrorsUnion `json:"errors,omitempty"`
	ProfileId *uint32      `json:"id,omitempty"`
}

func (response LoginResponse) Bind() {}
