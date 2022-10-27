package courses

import (
	"net/http"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
)

func collectCoursesByPeriod(context *dependency.DependencyContext, periodId uint32, userId uint32) (*CoursesByPeriod, error) {
	context.Logger.Debugf("List: start processing periodId = %d", periodId)

	coursesInDB, err := context.Storage.CourseDAO().GetCoursesByPeriod(periodId)
	if err != nil {
		context.Logger.Errorf("List: error in finding courses by periodId = %d", periodId)
		return nil, err
	}

	coursesInResponse := make([]*Course, len(coursesInDB))
	for index, courseInDB := range coursesInDB {
		var courseInResponse Course = Course{Id: courseInDB.Id, Title: courseInDB.Title,
			ShortDescription: courseInDB.ShortDescription}

		teachers, errTeachers := context.Storage.TeacherToCourseDAO().GetTeachersByCourse(courseInDB.Id)
		if errTeachers != nil {
			context.Logger.Errorf("List: error in finding teachers by course %d", courseInDB.Id)
			return nil, err
		}
		courseInResponse.Teachers = teachers

		enroll, errEnroll := context.Storage.StudentToCourseDAO().ExistRecord(userId, courseInDB.Id)
		if errEnroll != nil {
			context.Logger.Errorf("List: error in checking existence (%d, %d) in student_to_course_base", userId, courseInDB.Id)
			return nil, errEnroll
		}
		courseInResponse.Enroll = enroll

		coursesInResponse[index] = &courseInResponse
	}
	return &CoursesByPeriod{PeriodId: periodId, Courses: coursesInResponse}, nil
}

func HandleCoursesList(request *ListRequest, context *dependency.DependencyContext, extraParameters ...any) (int, *ListResponse) {
	if len(extraParameters) == 0 {
		context.Logger.Errorf("List: don't get login of current user")
		return http.StatusInternalServerError, &ListResponse{}
	}

	var login string
	login, ok := extraParameters[0].(string)
	if !ok {
		context.Logger.Errorf("List: can't parse login of current user")
		return http.StatusInternalServerError, &ListResponse{}
	}
	user, errUser := context.Storage.UserDAO().FindUserByLogin(login)
	if errUser != nil {
		context.Logger.Errorf("List: can't find user with login = %s, %s", login, errUser)
		return http.StatusBadRequest, &ListResponse{}
	}

	context.Logger.Debugf("List: start with periods = %v for user = %s", request.Periods, login)

	if len(request.Periods) == 0 {
		info, err := context.Storage.GeneralDAO().GetInfo()
		if err != nil {
			context.Logger.Errorf("List: can't find current semester, %s", err)
			return http.StatusInternalServerError, &ListResponse{}
		}
		request.Periods = append(request.Periods, info.CurrentPeriodId)
	}

	for _, periodId := range request.Periods {
		exist, errExist := context.Storage.PeriodDAO().ExistPeriod(periodId)
		if errExist != nil {
			context.Logger.Errorf("List: can't select periodId = %d in database period_base, %s", periodId, errExist)
			return http.StatusInternalServerError, &ListResponse{}
		}
		if !exist {
			context.Logger.Errorf("List: can't find periodId = %d", periodId)
			return http.StatusBadRequest, &ListResponse{Errors: &ErrorsUnion{InvalidPeriodId: &Error{}}}
		}
	}

	var coursesAndPeriods CoursesAndPeriods
	coursesAndPeriods.CoursesByPeriods = make([]*CoursesByPeriod, len(request.Periods))
	for index, periodId := range request.Periods {
		coursesByPeriod, errPeriod := collectCoursesByPeriod(context, periodId, user.ProfileId)
		if errPeriod != nil {
			context.Logger.Errorf("List: can't process period %d, %s", periodId, errPeriod)
			return http.StatusInternalServerError, &ListResponse{}
		}
		coursesAndPeriods.CoursesByPeriods[index] = coursesByPeriod
	}

	return http.StatusOK, &ListResponse{Response: &coursesAndPeriods}
}
