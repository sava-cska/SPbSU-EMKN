package profiles

import (
	"net/http"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
)

// CreateOrder godoc
// @Summary Get profile
// @Tags profiles
// @Accept  json
// @Produce  json
// @Param request body GetRequest true "Get profile"
// @Success 200 {object} GetResponse
// @Router /profiles/get [post]
func HandleProfilesGet(request *GetRequest, context *dependency.DependencyContext, _ ...any) (int, *GetResponse) {
	if request.ProfileIds == nil || len(request.ProfileIds) == 0 {
		return http.StatusOK, &GetResponse{}
	}
	users, err := context.Storage.UserDAO().FindUsers(request.ProfileIds)
	if err != nil {
		return http.StatusBadRequest, &GetResponse{}
	}
	profiles, err := context.Storage.UserAvatarDAO().GetProfileById(request.ProfileIds)
	if err != nil {
		return http.StatusBadRequest, &GetResponse{}
	}

	return http.StatusOK, &GetResponse{
		Response: &GetWrapper{
			Profiles: toResponse(users, profiles),
		},
	}
}

func toResponse(users []models.User, profiles []models.Profile) *[]Profile {
	var res []Profile
	for _, val := range users {
		res = append(res, Profile{
			Id:        val.ProfileId,
			FirstName: val.FirstName,
			LastName:  val.LastName,
			AvatarUrl: "",
		})
	}
	for _, val := range profiles {
		for idx, rs := range res {
			if rs.Id == val.ProfileId {
				res[idx].AvatarUrl = val.AvatarUrl
				break
			}
		}
	}
	return &res
}
