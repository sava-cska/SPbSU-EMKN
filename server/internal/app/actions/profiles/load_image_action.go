package profiles

import (
	"encoding/base64"
	"net/http"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/data_saver"
)

// CreateOrder godoc
// @Summary Load profile image
// @Tags profiles
// @Accept  json
// @Produce  json
// @Param request body LoadImageRequest true "Load profile image"
// @Success 200 {object} LoadImageResponse
// @Router /profiles/load_image [post]
func HandleProfilesLoadImage(request *LoadImageRequest, context *dependency.DependencyContext, args ...any) (int, *LoadImageResponse) {
	if v, ok := args[0].(string); ok {
		return handleProfilesLoadImage(request, context, v)
	} else {
		return http.StatusInternalServerError, &LoadImageResponse{}
	}
}

func handleProfilesLoadImage(request *LoadImageRequest, context *dependency.DependencyContext, login string) (int, *LoadImageResponse) {
	decodedJpg, err := base64.StdEncoding.DecodeString(request.EncodedJpg)
	if err != nil {
		return http.StatusBadRequest, &LoadImageResponse{}
	}

	var saver data_saver.DataSaver = &data_saver.JpgSaver{}
	fileName, err := saver.HardSave(decodedJpg)
	if err != nil {
		return http.StatusBadRequest, &LoadImageResponse{}
	}

	absoluteUrl := "http://51.250.98.212:8888/" + fileName

	user, _ := context.Storage.UserDAO().FindUserByLogin(login)

	err = context.Storage.UserAvatarDAO().UpdateProfile(models.Profile{
		ProfileId: user.ProfileId,
		AvatarUrl: absoluteUrl,
	})

	if err != nil {
		context.Logger.Errorf(err.Error())
		return http.StatusBadRequest, &LoadImageResponse{}
	}

	return http.StatusOK, &LoadImageResponse{
		Response: &LoadImageWrapper{
			Url: absoluteUrl,
		},
	}
}
