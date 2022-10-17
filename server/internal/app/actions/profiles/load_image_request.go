package profiles

type LoadImageRequest struct {
	EncodedJpg string `json:"encoded_jpg"`
}

func (request LoadImageRequest) Bind() {}
