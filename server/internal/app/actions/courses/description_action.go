package courses

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"net/http"
)

func HandleCoursesDescription(request *DescriptionRequest, context *dependency.DependencyContext, _ ...any) (int, *DescriptionResponse) {
	description, err := context.Storage.CoursesDAO().GetDescription(request.Id)
	if err != nil {
		context.Logger.Errorf("Failed to get course description for course id %d: %s", request.Id, err.Error())
		return http.StatusBadRequest, &DescriptionResponse{}
	}
	if description == nil {
		return http.StatusBadRequest, &DescriptionResponse{
			Errors: &ErrorsUnion{InvalidCourseId: &Error{}},
		}
	}

	return http.StatusOK, &DescriptionResponse{Response: &DescriptionResponseWrapper{ Description: *description}}
}