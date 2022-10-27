package courses

type Error struct{}

type ErrorsUnion struct {
	InvalidCourseId *Error `json:"invalid_course_id,omitempty"`
	InvalidPeriodId *Error `json:"invalid_period_id,omitempty"`
	AlreadyEnrolled *Error `json:"already_enrolled,omitempty"`
}
