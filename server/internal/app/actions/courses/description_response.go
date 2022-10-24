package courses

type DescriptionResponseWrapper struct {
	Description string `json:"description"`
}

type DescriptionResponse struct {
	Response *DescriptionResponseWrapper `json:"response,omitempty"`
	Errors   *ErrorsUnion                `json:"errors,omitempty"`
}

func (d DescriptionResponse) Bind() {}
