package courses

type UnEnrollResponse struct {
	Errors *ErrorsUnion `json:"errors,omitempty"`
}

func (response UnEnrollResponse) Bind() {}
