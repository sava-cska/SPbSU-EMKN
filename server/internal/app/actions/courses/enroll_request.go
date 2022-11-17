package courses

type EnrollRequest struct {
	CourseId uint32 `json:"course_id"`
}

func (request EnrollRequest) Bind() {}
