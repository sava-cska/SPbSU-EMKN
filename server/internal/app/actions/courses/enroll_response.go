package courses

type EnrollResponse struct {
	Errors *ErrorsUnion `json:"errors,omitempty"`
}

func (response EnrollResponse) Bind() {}
