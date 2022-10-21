package courses

type DescriptionRequest struct {
	Id uint `json:"id"`
}

func (d DescriptionRequest) Bind() {}
