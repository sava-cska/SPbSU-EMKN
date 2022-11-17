package courses

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage/test_storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sort"
	"testing"
	"time"
)

func TestDescription(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		exp := time.Now()
		cont, db, _ := dependency.NewTestContext()
		db.AddNewCourse(test_storage.TestCourseData{}.Fill(0, "Math", "Mathematics", &exp, 0, make([]uint32, 0), make([]uint32, 0)))
		code, resp := HandleCoursesDescription(&DescriptionRequest{Id: 0}, cont)

		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, "Mathematics", resp.Response.Description)
	})

	t.Run("Invalid course", func(t *testing.T) {
		exp := time.Now()
		cont, db, _ := dependency.NewTestContext()
		db.AddNewCourse(test_storage.TestCourseData{}.Fill(1, "Math", "Mathematics", &exp, 0, make([]uint32, 0), make([]uint32, 0)))
		code, resp := HandleCoursesDescription(&DescriptionRequest{Id: 0}, cont)

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, ErrorsUnion{
			InvalidCourseId: &Error{},
		}, *resp.Errors)
	})
}

func TestDescriptionping(t *testing.T) {
	exp := time.Now()
	cont, db, _ := dependency.NewTestContext()
	db.AddNewCourse(test_storage.TestCourseData{}.Fill(0, "Math", "Mathematics", &exp, 0, make([]uint32, 0), make([]uint32, 0)))

	t.Run("Successful", func(t *testing.T) {

		lastSyncTime := time.Now().Add(1 * time.Second)
		code, resp := HandleCoursesDescriptionPing(&DescriptionPingRequest{Id: 0, LastSyncTime: lastSyncTime.UnixMilli()}, cont)

		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, false, *resp.Response.Changed)

		lastSyncTime = time.Now().Add(-1 * time.Second)
		code, resp = HandleCoursesDescriptionPing(&DescriptionPingRequest{Id: 0, LastSyncTime: lastSyncTime.UnixMilli()}, cont)

		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, true, *resp.Response.Changed)
	})

	t.Run("Invalid course id", func(t *testing.T) {

		lastSyncTime := time.Now().Add(1 * time.Second)
		code, resp := HandleCoursesDescriptionPing(&DescriptionPingRequest{Id: 1, LastSyncTime: lastSyncTime.UnixMilli()}, cont)

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, ErrorsUnion{
			InvalidCourseId: &Error{},
		}, *resp.Errors)
	})
}

func TestEnroll(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		exp := time.Now()
		cont, db, _ := dependency.NewTestContext()

		db.AddNewCourse(test_storage.TestCourseData{}.Fill(0, "Math", "Mathematics", &exp, 0, make([]uint32, 0), make([]uint32, 0)))
		db.AddNewUser(1, "jane_doe", "qwerty", "jane.doe@gmail.com", "Jane", "Doe")

		code, _ := HandleCoursesEnroll(&EnrollRequest{CourseId: 0}, cont, "jane_doe")

		assert.Equal(t, http.StatusOK, code)

		course := db.Courses[0]
		assert.Equal(t, true, course.Students[1])
		assert.Equal(t, 1, len(course.Students))
	})

	t.Run("Already enrolled", func(t *testing.T) {
		exp := time.Now()
		cont, db, _ := dependency.NewTestContext()

		db.AddNewCourse(test_storage.TestCourseData{}.Fill(0, "Math", "Mathematics", &exp, 0, []uint32{1}, make([]uint32, 0)))
		db.AddNewUser(1, "jane_doe", "qwerty", "jane.doe@gmail.com", "Jane", "Doe")

		code, resp := HandleCoursesEnroll(&EnrollRequest{CourseId: 0}, cont, "jane_doe")

		course := db.Courses[0]
		assert.Equal(t, true, course.Students[1])
		assert.Equal(t, 1, len(course.Students))

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, ErrorsUnion{
			AlreadyEnrolled: &Error{},
		}, *resp.Errors)
	})
}

func TestList(t *testing.T) {
	exp := time.Now()
	cont, db, _ := dependency.NewTestContext()

	db.AddNewUser(1, "jane_doe", "qwerty", "jane.doe@gmail.com", "Jane", "Doe")
	db.AddNewUser(2, "norma_jean", "asdfasdf", "norma.jean@gmail.com", "Norma", "Jean")

	db.AddNewCourse(test_storage.TestCourseData{}.Fill(0, "Math", "Mathematics", &exp, 0, []uint32{1, 2}, make([]uint32, 0)))
	db.AddNewCourse(test_storage.TestCourseData{}.Fill(1, "Ph", "Philosophy", &exp, 0, []uint32{1}, make([]uint32, 0)))

	db.AddNewCourse(test_storage.TestCourseData{}.Fill(2, "Math2", "Mathematics", &exp, 1, []uint32{2}, make([]uint32, 0)))
	db.AddNewCourse(test_storage.TestCourseData{}.Fill(3, "Ph2", "Philosophy", &exp, 1, []uint32{2, 1}, make([]uint32, 0)))

	db.AddPeriod(0, "2022")
	db.AddPeriod(1, "2023")
	db.AddPeriod(3, "2024")

	db.CurrentPeriodId = 1

	t.Run("Successful", func(t *testing.T) {
		code, resp := HandleCoursesList(&ListRequest{Periods: []uint32{0, 1}}, cont, "jane_doe")

		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, 2, len(resp.Response.CoursesByPeriods))
		sort.Slice(resp.Response.CoursesByPeriods[0].Courses, func(i, j int) bool {
			return resp.Response.CoursesByPeriods[0].Courses[i].Id < resp.Response.CoursesByPeriods[0].Courses[j].Id
		})
		sort.Slice(resp.Response.CoursesByPeriods[1].Courses, func(i, j int) bool {
			return resp.Response.CoursesByPeriods[1].Courses[i].Id < resp.Response.CoursesByPeriods[1].Courses[j].Id
		})

		assert.Equal(t, uint32(0), resp.Response.CoursesByPeriods[0].PeriodId)
		assert.Equal(t, 2, len(resp.Response.CoursesByPeriods[0].Courses))
		assert.Equal(t, uint32(0), resp.Response.CoursesByPeriods[0].Courses[0].Id)
		assert.True(t, resp.Response.CoursesByPeriods[0].Courses[0].Enroll)
		assert.Equal(t, uint32(1), resp.Response.CoursesByPeriods[0].Courses[1].Id)
		assert.True(t, resp.Response.CoursesByPeriods[0].Courses[1].Enroll)

		assert.Equal(t, uint32(1), resp.Response.CoursesByPeriods[1].PeriodId)
		assert.Equal(t, 2, len(resp.Response.CoursesByPeriods[1].Courses))
		assert.Equal(t, uint32(2), resp.Response.CoursesByPeriods[1].Courses[0].Id)
		assert.False(t, resp.Response.CoursesByPeriods[1].Courses[0].Enroll)
		assert.Equal(t, uint32(3), resp.Response.CoursesByPeriods[1].Courses[1].Id)
		assert.True(t, resp.Response.CoursesByPeriods[1].Courses[1].Enroll)
	})

	t.Run("Current period", func(t *testing.T) {
		code, resp := HandleCoursesList(&ListRequest{Periods: []uint32{}}, cont, "jane_doe")

		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, 1, len(resp.Response.CoursesByPeriods))
		sort.Slice(resp.Response.CoursesByPeriods[0].Courses, func(i, j int) bool {
			return resp.Response.CoursesByPeriods[0].Courses[i].Id < resp.Response.CoursesByPeriods[0].Courses[j].Id
		})
		assert.Equal(t, uint32(1), resp.Response.CoursesByPeriods[0].PeriodId)
		assert.Equal(t, 2, len(resp.Response.CoursesByPeriods[0].Courses))
		assert.Equal(t, uint32(2), resp.Response.CoursesByPeriods[0].Courses[0].Id)
		assert.False(t, resp.Response.CoursesByPeriods[0].Courses[0].Enroll)
		assert.Equal(t, uint32(3), resp.Response.CoursesByPeriods[0].Courses[1].Id)
		assert.True(t, resp.Response.CoursesByPeriods[0].Courses[1].Enroll)
	})

	t.Run("Invalid period", func(t *testing.T) {

		code, resp := HandleCoursesList(&ListRequest{Periods: []uint32{1, 4}}, cont, "jane_doe")

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, ErrorsUnion{
			InvalidPeriodId: &Error{},
		}, *resp.Errors)
	})
}

func TestPeriods(t *testing.T) {
	cont, db, _ := dependency.NewTestContext()

	db.AddPeriod(0, "2022")
	db.AddPeriod(1, "2023")
	db.AddPeriod(2, "2024")

	db.CurrentPeriodId = 1

	t.Run("Successful", func(t *testing.T) {

		code, resp := HandleCoursesPeriods(&PeriodsRequest{}, cont)

		assert.Equal(t, http.StatusOK, code)

		assert.Equal(t, 3, len(resp.Response.Periods))
		assert.Equal(t, uint32(1), resp.Response.Periods[0].Id)
		sort.Slice(resp.Response.Periods, func(i, j int) bool {
			return resp.Response.Periods[i].Id < resp.Response.Periods[j].Id
		})
		assert.Equal(t, []uint32{0, 1, 2}, []uint32{resp.Response.Periods[0].Id, resp.Response.Periods[1].Id, resp.Response.Periods[2].Id})
	})
}

func TestUnenroll(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		exp := time.Now()
		cont, db, _ := dependency.NewTestContext()

		db.AddNewCourse(test_storage.TestCourseData{}.Fill(0, "Math", "Mathematics", &exp, 0, []uint32{2, 1}, make([]uint32, 0)))
		db.AddNewUser(1, "jane_doe", "qwerty", "jane.doe@gmail.com", "Jane", "Doe")
		db.AddNewUser(2, "norma_jean", "asdfasdf", "norma.jean@gmail.com", "Norma", "Jean")

		code, _ := HandleCoursesUnEnroll(&UnEnrollRequest{CourseId: 0}, cont, "jane_doe")

		assert.Equal(t, http.StatusOK, code)

		course := db.Courses[0]
		assert.Equal(t, true, course.Students[2])
		assert.Equal(t, 1, len(course.Students))
	})

	t.Run("Not enrolled", func(t *testing.T) {
		exp := time.Now()
		cont, db, _ := dependency.NewTestContext()

		db.AddNewCourse(test_storage.TestCourseData{}.Fill(0, "Math", "Mathematics", &exp, 0, []uint32{2}, make([]uint32, 0)))
		db.AddNewUser(1, "jane_doe", "qwerty", "jane.doe@gmail.com", "Jane", "Doe")
		db.AddNewUser(2, "norma_jean", "asdfasdf", "norma.jean@gmail.com", "Norma", "Jean")

		code, resp := HandleCoursesUnEnroll(&UnEnrollRequest{CourseId: 0}, cont, "jane_doe")

		course := db.Courses[0]
		assert.Equal(t, true, course.Students[2])
		assert.Equal(t, 1, len(course.Students))

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, ErrorsUnion{
			CourseIsNotEnrolled: &Error{},
		}, *resp.Errors)
	})
}
