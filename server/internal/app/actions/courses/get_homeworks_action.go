package courses

import (
	"net/http"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
)

// CreateOrder godoc
// @Summary GetHomeworks course
// @Tags courses
// @Accept  json
// @Produce  json
// @Param request body GetHomeworksRequest true "GetHomeworks course"
// @Success 200 {object} GetHomeworksResponse
// @Router /courses/get_homeworks [post]
func HandleCoursesGetHomeworks(request *GetHomeworksRequest, context *dependency.DependencyContext,
	extraParameters ...any) (int, *GetHomeworksResponse) {
	if len(extraParameters) == 0 {
		context.Logger.Errorf("GetHomeworks: don't get login of current user")
		return http.StatusInternalServerError, &GetHomeworksResponse{}
	}

	var login string
	login, ok := extraParameters[0].(string)
	if !ok {
		context.Logger.Errorf("GetHomeworks: can't parse login of current user")
		return http.StatusInternalServerError, &GetHomeworksResponse{}
	}
	user, errUser := context.Storage.UserDAO().FindUserByLogin(login)
	if errUser != nil {
		context.Logger.Errorf("GetHomeworks: can't find user with login = %s, %s", login, errUser)
		return http.StatusInternalServerError, &GetHomeworksResponse{}
	}

	context.Logger.Debugf("GetHomeworks: start with courseId = %d for user = %s", request.CourseId, login)

	exist, errExist := context.Storage.CourseDAO().ExistCourse(request.CourseId)
	if errExist != nil {
		context.Logger.Errorf("GetHomeworks: error in checking existence %d in course_base, %s", request.CourseId, errExist)
		return http.StatusInternalServerError, &GetHomeworksResponse{}
	}
	if !exist {
		context.Logger.Errorf("GetHomeworks: invalid courseId %d", request.CourseId)
		return http.StatusBadRequest, &GetHomeworksResponse{Errors: &ErrorsUnion{InvalidCourseId: &Error{}}}
	}

	homeworks, errAllHomeworks := context.Storage.HomeworkDAO().GetAllHomeworks(request.CourseId)
	if errAllHomeworks != nil {
		context.Logger.Errorf("GetHomeworks: can't get all homeworks for course %d from homework_base, %s", request.CourseId,
			errAllHomeworks)
		return http.StatusInternalServerError, &GetHomeworksResponse{}
	}

	var response_homeworks []Homework
	for _, homework := range homeworks {
		response_homework := Homework{Id: homework.Id, Name: homework.Name, Deadline: homework.Deadline}

		checked, score, errChecked := context.Storage.CheckedHomeworkDAO().ScoreForHomework(user.ProfileId, homework.Id)
		if errChecked != nil {
			context.Logger.Errorf("GetHomeworks: can't get score for homework %d from checked_homework_base, %s", homework.Id,
				errChecked)
			continue
		}
		if checked {
			response_homework.StatusChecked = &HomeworkStatusChecked{TotalScore: homework.TotalScore, Score: score}
			response_homeworks = append(response_homeworks, response_homework)
			continue
		}

		passed, errPassed := context.Storage.PassedHomeworkDAO().CheckUserPassHomework(user.ProfileId, homework.Id)
		if errPassed != nil {
			context.Logger.Errorf("GetHomeworks: can't get information about homework %d from passed_homework_base, %s",
				homework.Id, errPassed)
			continue
		}
		if passed {
			response_homework.StatusUnchecked = &HomeworkStatusUnchecked{Message: "Не проверено"}
		} else {
			response_homework.StatusNotPassed = &HomeworkStatusNotPassed{Message: "Не сдано"}
		}
		response_homeworks = append(response_homeworks, response_homework)
	}

	return http.StatusOK, &GetHomeworksResponse{Response: &HomeworksList{Homeworks: &response_homeworks}}
}
