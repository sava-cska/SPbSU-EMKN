package courses

type DescriptionPingResponseWrapper struct {
	Changed *bool        `json:"changed,omitempty"`
}

type DescriptionPingResponse struct {
	Response *DescriptionPingResponseWrapper `json:"response,omitempty"`
	Errors  *ErrorsUnion `json:"errors,omitempty"`
}

func (d DescriptionPingResponse) Bind() {}
