package courses

import (
	"net/http"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
)

func HandleCoursesUnEnroll(request *UnEnrollRequest, context *dependency.DependencyContext,
	extraParameters ...any) (int, *UnEnrollResponse) {
	if len(extraParameters) == 0 {
		context.Logger.Errorf("UnEnroll: don't get login of current user")
		return http.StatusInternalServerError, &UnEnrollResponse{}
	}

	var login string
	login, ok := extraParameters[0].(string)
	if !ok {
		context.Logger.Errorf("UnEnroll: can't parse login of current user")
		return http.StatusInternalServerError, &UnEnrollResponse{}
	}
	user, errUser := context.Storage.UserDAO().FindUserByLogin(login)
	if errUser != nil {
		context.Logger.Errorf("UnEnroll: can't find user with login = %s, %s", login, errUser)
		return http.StatusInternalServerError, &UnEnrollResponse{}
	}

	context.Logger.Debugf("UnEnroll: start with courseId = %d for user = %s", request.CourseId, login)

	exist, errExist := context.Storage.CourseDAO().ExistCourse(request.CourseId)
	if errExist != nil {
		context.Logger.Errorf("UnEnroll: error in checking existence %d in course_base", request.CourseId)
		return http.StatusInternalServerError, &UnEnrollResponse{}
	}
	if !exist {
		context.Logger.Errorf("UnEnroll: invalid courseId %d", request.CourseId)
		return http.StatusBadRequest, &UnEnrollResponse{Errors: &ErrorsUnion{InvalidCourseId: &Error{}}}
	}

	exist, err := context.Storage.StudentToCourseDAO().ExistRecord(user.ProfileId, request.CourseId)
	if err != nil {
		context.Logger.Errorf("UnEnroll: error in checking existence (%d, %d) in student_to_course_base",
			user.ProfileId, request.CourseId)
		return http.StatusInternalServerError, &UnEnrollResponse{}
	}
	if !exist {
		context.Logger.Errorf("UnEnroll: user %s wasn't enrolled to course %d", login, request.CourseId)
		return http.StatusBadRequest, &UnEnrollResponse{Errors: &ErrorsUnion{CourseIsNotEnrolled: &Error{}}}
	}

	if errDelete := context.Storage.StudentToCourseDAO().DeleteRecord(user.ProfileId, request.CourseId); errDelete != nil {
		context.Logger.Errorf("UnEnroll: error in deleting (%d, %d) from student_to_course_base", user.ProfileId, request.CourseId)
		return http.StatusInternalServerError, &UnEnrollResponse{}
	}

	return http.StatusOK, &UnEnrollResponse{}
}
