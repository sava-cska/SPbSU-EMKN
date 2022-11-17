package courses

type Course struct {
	Id               uint32   `json:"id"`
	Title            string   `json:"title"`
	Enroll           bool     `json:"enrolled"`
	ShortDescription string   `json:"short_description"`
	Teachers         []uint32 `json:"teachers_profiles"`
}

type CoursesByPeriod struct {
	PeriodId uint32    `json:"period_id"`
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
