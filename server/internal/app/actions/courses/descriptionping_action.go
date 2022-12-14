package courses

import (
	"net/http"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
)

// CreateOrder godoc
// @Summary Ping course description
// @Tags courses
// @Accept  json
// @Produce  json
// @Param request body DescriptionPingRequest true "Ping course description"
// @Success 200 {object} DescriptionPingResponse
// @Router /courses/description_ping [post]
func HandleCoursesDescriptionPing(request *DescriptionPingRequest, context *dependency.DependencyContext, _ ...any) (int, *DescriptionPingResponse) {
	timestamp, err := context.Storage.CourseDAO().GetDescriptionTimestamp(request.Id)
	if err != nil {
		context.Logger.Errorf("Failed to get course description timestamp for course id %d: %s", request.Id, err.Error())
		return http.StatusBadRequest, &DescriptionPingResponse{}
	}
	if timestamp == nil {
		return http.StatusBadRequest, &DescriptionPingResponse{
			Errors: &ErrorsUnion{InvalidCourseId: &Error{}},
		}
	}

	context.Logger.Debugf("CoursesDescriptionPing: Last change timestamp: %d, timestamp from reqeust %d", timestamp.UnixMilli(), request.LastSyncTime)
	changed := timestamp.UnixMilli() > request.LastSyncTime
	return http.StatusOK, &DescriptionPingResponse{Response: &DescriptionPingResponseWrapper{Changed: &changed}}
}
