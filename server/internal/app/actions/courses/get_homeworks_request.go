package courses

type GetHomeworksRequest struct {
	CourseId uint32 `json:"course_id"`
}

func (request GetHomeworksRequest) Bind() {}
