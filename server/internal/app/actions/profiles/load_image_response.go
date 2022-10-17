package profiles

type LoadImageWrapper struct {
	Url string `json:"url"`
}

type LoadImageResponse struct {
	Errors   *ErrorsUnion      `json:"errors,omitempty"`
	Response *LoadImageWrapper `json:"response,omitempty"`
}

func (response LoadImageResponse) Bind() {}
