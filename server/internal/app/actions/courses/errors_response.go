package courses

type Error struct{}

type ErrorsUnion struct {
	InvalidCourseId     *Error `json:"invalid_course_id,omitempty"`
	InvalidPeriodId     *Error `json:"invalid_period_id,omitempty"`
	AlreadyEnrolled     *Error `json:"already_enrolled,omitempty"`
	CourseIsNotEnrolled *Error `json:"course_is_not_enrolled,omitempty"`
}
