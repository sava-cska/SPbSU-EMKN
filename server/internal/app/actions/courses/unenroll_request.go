package courses

type UnEnrollRequest struct {
	CourseId uint32 `json:"course_id"`
}

func (request UnEnrollRequest) Bind() {}
