package core

import (
	"encoding/json"
	"fmt"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage/test_storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestRegistration(t *testing.T) {
	server, db, mailer := NewTest()
	ok, id := registerUser(&models.User{
		Login:     "jane_doe",
		Password:  "qwerty",
		Email:     "jane_doe@gmail.com",
		FirstName: "Jane",
		LastName:  "Doe",
	}, server, db, mailer)
	if !ok {
		t.Fatalf("Failed to create user")
	}

	user := db.LoginToUser["jane_doe"]
	assert.Equal(t, "jane_doe", user.Login)
	assert.Equal(t, "qwerty", user.Password)
	assert.Equal(t, id, user.ProfileId)
	assert.Equal(t, "Jane", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "jane_doe@gmail.com", user.Email)
}

func registerUser(user *models.User, server *TestServer, db *test_storage.TestDAO, mailer *dependency.TestMailer) (bool, uint32) {
	statusCode, regResp := server.SendRequest("/accounts/register", fmt.Sprintf(`
    {
      "login": "%s",
      "password": "%s",
      "email": "%s",
      "first_name": "%s",
      "last_name": "%s"
    }`, user.Login, user.Password, user.Email, user.FirstName, user.LastName))

	if statusCode != http.StatusOK {
		return false, 0
	}

	mailer.WaitSingleMessageSent()
	msg := mailer.Msgs["jane_doe@gmail.com"][0]
	verificationCode := parseVerificationCode(msg)

	respMap := unmarshallMap(regResp)
	respMap = respMap["response"].(map[string]interface{})
	randomToken := respMap["random_token"].(string)

	statusCode, _ = server.SendRequest("/accounts/validate_email", fmt.Sprintf(`
		{
			"verification_code": "%s",
			"random_token": "%s"
		}`, verificationCode, randomToken))

	if statusCode != http.StatusOK {
		return false, 0
	}

	addedUser, err := db.FindUserByLogin(user.Login)
	if err != nil {
		return false, 0
	}
	return true, addedUser.ProfileId
}

func parseVerificationCode(msg *notifier.Message) string {
	subs := "Код подтверждения: <b>"
	ind := strings.Index(msg.Body, subs) + len(subs)
	verificationCode := msg.Body[ind : ind+6]
	return verificationCode
}

func unmarshallMap(jsonStr string) map[string]interface{} {
	res := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &res)
	if err != nil {
		panic(err)
	}
	return res
}
