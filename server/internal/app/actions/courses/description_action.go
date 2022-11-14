package courses

import (
	"net/http"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
)

// CreateOrder godoc
// @Summary Get course description
// @Tags courses
// @Accept  json
// @Produce  json
// @Param request body DescriptionRequest true "Get course description"
// @Success 200 {object} DescriptionResponse
// @Router /courses/description [post]
func HandleCoursesDescription(request *DescriptionRequest, context *dependency.DependencyContext, _ ...any) (int, *DescriptionResponse) {
	description, err := context.Storage.CourseDAO().GetDescription(request.Id)
	if err != nil {
		context.Logger.Errorf("Failed to get course description for course id %d: %s", request.Id, err.Error())
		return http.StatusBadRequest, &DescriptionResponse{}
	}
	if description == nil {
		return http.StatusBadRequest, &DescriptionResponse{
			Errors: &ErrorsUnion{InvalidCourseId: &Error{}},
		}
	}

	return http.StatusOK, &DescriptionResponse{Response: &DescriptionResponseWrapper{Description: *description}}
}
