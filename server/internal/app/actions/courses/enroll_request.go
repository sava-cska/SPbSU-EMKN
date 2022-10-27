package courses

type EnrollRequest struct {
	CourseId uint32 `json:"course_id,omitempty"`
}

func (request EnrollRequest) Bind() {}
