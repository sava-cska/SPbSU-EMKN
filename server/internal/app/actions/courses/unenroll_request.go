package courses

type UnEnrollRequest struct {
	CourseId uint32 `json:"course_id,omitempty"`
}

func (request UnEnrollRequest) Bind() {}
