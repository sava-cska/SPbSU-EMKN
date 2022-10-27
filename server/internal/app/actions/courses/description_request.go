package courses

type DescriptionRequest struct {
	Id uint32 `json:"id"`
}

func (d DescriptionRequest) Bind() {}
