package courses

import (
	"net/http"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
)

// CreateOrder godoc
// @Summary Enroll course
// @Tags courses
// @Accept  json
// @Produce  json
// @Param request body EnrollRequest true "Enroll course"
// @Success 200 {object} EnrollResponse
// @Router /courses/enroll [post]
func HandleCoursesEnroll(request *EnrollRequest, context *dependency.DependencyContext, extraParameters ...any) (int, *EnrollResponse) {
	if len(extraParameters) == 0 {
		context.Logger.Errorf("Enroll: don't get login of current user")
		return http.StatusInternalServerError, &EnrollResponse{}
	}

	var login string
	login, ok := extraParameters[0].(string)
	if !ok {
		context.Logger.Errorf("Enroll: can't parse login of current user")
		return http.StatusInternalServerError, &EnrollResponse{}
	}
	user, errUser := context.Storage.UserDAO().FindUserByLogin(login)
	if errUser != nil {
		context.Logger.Errorf("Enroll: can't find user with login = %s, %s", login, errUser)
		return http.StatusInternalServerError, &EnrollResponse{}
	}

	context.Logger.Debugf("Enroll: start with courseId = %d for user = %s", request.CourseId, login)

	exist, errExist := context.Storage.CourseDAO().ExistCourse(request.CourseId)
	if errExist != nil {
		context.Logger.Errorf("Enroll: error in checking existence %d in course_base, %s", request.CourseId, errExist)
		return http.StatusInternalServerError, &EnrollResponse{}
	}
	if !exist {
		context.Logger.Errorf("Enroll: invalid courseId %d", request.CourseId)
		return http.StatusBadRequest, &EnrollResponse{Errors: &ErrorsUnion{InvalidCourseId: &Error{}}}
	}

	exist, err := context.Storage.StudentToCourseDAO().ExistRecord(user.ProfileId, request.CourseId)
	if err != nil {
		context.Logger.Errorf("Enroll: error in checking existence (%d, %d) in student_to_course_base, %s",
			user.ProfileId, request.CourseId, err)
		return http.StatusInternalServerError, &EnrollResponse{}
	}
	if exist {
		context.Logger.Errorf("Enroll: user %s already enrolled to course %d", login, request.CourseId)
		return http.StatusBadRequest, &EnrollResponse{Errors: &ErrorsUnion{AlreadyEnrolled: &Error{}}}
	}

	if errAdd := context.Storage.StudentToCourseDAO().AddRecord(user.ProfileId, request.CourseId); errAdd != nil {
		context.Logger.Errorf("Enroll: error in adding (%d, %d) into student_to_course_base, %s", user.ProfileId, request.CourseId, errAdd)
		return http.StatusInternalServerError, &EnrollResponse{}
	}

	return http.StatusOK, &EnrollResponse{}
}
