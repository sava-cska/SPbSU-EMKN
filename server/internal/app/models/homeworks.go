package models

type HomeworkInDB struct {
	Id         uint32
	Name       string
	Deadline   int64
	TotalScore int
	CourseId   uint32
}
