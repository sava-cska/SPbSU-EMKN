package courses

type Course struct {
	Id               uint32   `json:"id,omitempty"`
	Title            string   `json:"title,omitempty"`
	Enroll           bool     `json:"enrolled,omitempty"`
	ShortDescription string   `json:"short_description,omitempty"`
	Teachers         []uint32 `json:"teachers_profiles,omitempty"`
}

type CoursesByPeriod struct {
	PeriodId uint32    `json:"period_id,omitempty"`
	Courses  []*Course `json:"courses,omitempty"`
}

type CoursesAndPeriods struct {
	CoursesByPeriods []*CoursesByPeriod `json:"courses_by_period,omitempty"`
}

type ListResponse struct {
	Response *CoursesAndPeriods `json:"response,omitempty"`
	Errors   *ErrorsUnion       `json:"errors,omitempty"`
}

func (response ListResponse) Bind() {}
