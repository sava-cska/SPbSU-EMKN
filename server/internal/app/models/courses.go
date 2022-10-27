package models

type Period struct {
	Id   uint32 `json:"id"`
	Text string `json:"text"`
}

type CourseInDB struct {
	Id               uint32
	Title            string
	ShortDescription string
}
