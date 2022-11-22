package test_storage

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"time"
)

type TestRegistrationData struct {
	Token            string
	Login            string
	Password         string
	Email            string
	FirstName        string
	LastName         string
	ExpireData       *time.Time
	VerificationCode string
}

type TestCourseData struct {
	Course        *models.CourseInDB
	DescTimestamp *time.Time
	PeriodId      uint32
	Students      map[uint32]bool
	Teachers      map[uint32]bool
}

func (t TestCourseData) Fill(id uint32, title string, descr string, descrTime *time.Time, periodId uint32,
	studentIds []uint32, teacherIds []uint32) *TestCourseData {
	t.Students = make(map[uint32]bool)
	t.Teachers = make(map[uint32]bool)
	for _, st := range studentIds {
		t.Students[st] = true
	}
	for _, te := range teacherIds {
		t.Teachers[te] = true
	}
	t.Course = &models.CourseInDB{
		Id:               id,
		Title:            title,
		ShortDescription: descr,
	}
	t.DescTimestamp = descrTime
	t.PeriodId = periodId
	return &t
}

type TestChangePasswordData struct {
	Token                    string
	Login                    string
	ExpireDate               *time.Time
	VerificationCode         string
	ChangePasswordToken      string
	ChangePasswordExpireTime *time.Time
}

func (t TestChangePasswordData) Fill(token string, login string, expireDate *time.Time, verificationCode string) *TestChangePasswordData {
	t.Token = token
	t.ExpireDate = expireDate
	t.Login = login
	t.VerificationCode = verificationCode
	t.ChangePasswordExpireTime = &time.Time{}
	return &t
}
