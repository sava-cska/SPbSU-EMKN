package courses

type ListRequest struct {
	Periods []uint32 `json:"period_ids"`
}

func (request ListRequest) Bind() {}
