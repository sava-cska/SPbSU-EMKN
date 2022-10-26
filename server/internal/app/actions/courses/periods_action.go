package courses

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"net/http"
)

func HandleCoursesPeriods(request *PeriodsRequest, context *dependency.DependencyContext, _ ...any) (int, *PeriodsResponse) {
	periods, err := context.Storage.CoursesDAO().GetPeriods()
	if err != nil {
		context.Logger.Errorf("Failed to get periods from db: %s", err.Error())
		return http.StatusBadRequest, &PeriodsResponse{}
	}

	info, err := context.Storage.GeneralDAO().GetInfo()
	if err != nil {
		context.Logger.Errorf("Failed to get periods from db: %s", err.Error())
		return http.StatusInternalServerError, &PeriodsResponse{}
	}
	for i := 0; i < len(periods); i++ {
		if periods[i].Id == info.CurrentPeriodId {
			periods[0], periods[i] = periods[i], periods[0]
			break
		}
	}

	return http.StatusOK, &PeriodsResponse{&PeriodsResponseWrapper{Periods: periods}}
}
