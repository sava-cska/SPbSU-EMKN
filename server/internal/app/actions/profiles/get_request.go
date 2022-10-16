package profiles

type GetRequest struct {
	ProfileIds []int32 `json:"profile_ids"`
}

func (request GetRequest) Bind() {}
