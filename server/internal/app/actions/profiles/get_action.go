package profiles

import (
	"net/http"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
)

func HandleProfilesGet(request *GetRequest, context *dependency.DependencyContext, _ ...any) (int, *GetResponse) {
	if request.ProfileIds == nil || len(request.ProfileIds) == 0 {
		return http.StatusOK, &GetResponse{}
	}
	profiles, err := context.Storage.UserAvatarDAO().GetProfileById(request.ProfileIds)
	if err != nil {
		return http.StatusBadRequest, &GetResponse{}
	}

	return http.StatusOK, &GetResponse{
		Response: &GetWrapper{
			Profiles: toResponse(profiles),
		},
	}
}

func toResponse(profiles []models.Profile) *[]Profile {
	var res []Profile
	for _, val := range profiles {
		res = append(res, Profile{
			Id:        val.ProfileId,
			AvatarUrl: val.AvatarUrl,
			FirstName: val.FirstName,
			LastName:  val.LastName,
		})
	}
	return &res
}
