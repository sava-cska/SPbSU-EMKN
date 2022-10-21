package courses

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"net/http"
)

func HandleCoursesDescriptionPing(request *DescriptionPingRequest, context *dependency.DependencyContext, _ ...any) (int, *DescriptionPingResponse) {
	timestamp, err := context.Storage.CoursesDAO().GetDescriptionTimestamp(request.Id)
	if err != nil {
		context.Logger.Errorf("Failed to get course description timestamp for course id %d: %s", request.Id, err.Error())
		return http.StatusBadRequest, &DescriptionPingResponse{}
	}
	if timestamp == nil {
		return http.StatusBadRequest, &DescriptionPingResponse{
			Errors: &ErrorsUnion{InvalidCourseId: &Error{}},
		}
	}

	println("%d %d", timestamp.UnixMilli(), request.LastSyncTime)
	changed := timestamp.UnixMilli() > request.LastSyncTime
	return http.StatusOK, &DescriptionPingResponse{Response: &DescriptionPingResponseWrapper { Changed: &changed} }
}
