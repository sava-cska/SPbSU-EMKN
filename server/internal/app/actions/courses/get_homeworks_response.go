package courses

type HomeworkStatusNotPassed struct {
	Message string `json:"text"`
}

type HomeworkStatusUnchecked struct {
	Message string `json:"text"`
}

type HomeworkStatusChecked struct {
	TotalScore int `json:"total_score"`
	Score      int `json:"score"`
}

type Homework struct {
	Id              uint32                   `json:"id"`
	Name            string                   `json:"name"`
	Deadline        int64                    `json:"deadline"`
	StatusNotPassed *HomeworkStatusNotPassed `json:"status_not_passed,omitempty"`
	StatusUnchecked *HomeworkStatusUnchecked `json:"status_unchecked,omitempty"`
	StatusChecked   *HomeworkStatusChecked   `json:"status_checked,omitempty"`
}

type HomeworksList struct {
	Homeworks *[]Homework `json:"homeworks,omitempty"`
}

type GetHomeworksResponse struct {
	Response *HomeworksList `json:"response,omitempty"`
	Errors   *ErrorsUnion   `json:"errors,omitempty"`
}

func (response GetHomeworksResponse) Bind() {}
