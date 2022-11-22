package profiles

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sort"
	"testing"
)

func TestGet(t *testing.T) {
	cont, db, _ := dependency.NewTestContext()

	db.AddNewUser(1, "jane_doe", "qwerty", "jane.doe@gmail.com", "Jane", "Doe")
	db.AddNewUser(2, "norma_jean", "asdfasdf", "norma.jean@gmail.com", "Norma", "Jean")
	db.AddUserProfile(1, "www.google.com/jane_doe.jpg", "Jane", "Doe")
	db.AddUserProfile(2, "www.google.com/norma_jean.jpg", "Jane", "Doe")

	t.Run("Successful", func(t *testing.T) {
		code, resp := HandleProfilesGet(&GetRequest{ProfileIds: []int32{1, 2}}, cont)

		assert.Equal(t, http.StatusOK, code)
		sort.Slice(*resp.Response.Profiles, func(i, j int) bool {
			return (*resp.Response.Profiles)[i].Id < (*resp.Response.Profiles)[j].Id
		})
		assert.Equal(t, "www.google.com/jane_doe.jpg", (*resp.Response.Profiles)[0].AvatarUrl)
		assert.Equal(t, "www.google.com/norma_jean.jpg", (*resp.Response.Profiles)[1].AvatarUrl)
	})

	t.Run("Invalid profile id", func(t *testing.T) {
		code, _ := HandleProfilesGet(&GetRequest{ProfileIds: []int32{1, 2, 3}}, cont)

		assert.Equal(t, http.StatusBadRequest, code)
	})
}
