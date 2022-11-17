package actions

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/actions/accounts"
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
	cont, db, mailer := dependency.NewTestContext()
	ok, id := registerUser(&models.User{
		Login:     "jane_doe",
		Password:  "qwerty",
		Email:     "jane_doe@gmail.com",
		FirstName: "Jane",
		LastName:  "Doe",
	}, cont, db, mailer)
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

func TestParallelSameLoginRegistration(t *testing.T) {
	cont, db, mailer := dependency.NewTestContext()

	statusCode, resp1 := accounts.HandleAccountsRegister(&accounts.RegisterRequest{
		Login:     "jane_doe",
		Password:  "qwerty",
		Email:     "jane_doe@gmail.com",
		FirstName: "Jane",
		LastName:  "Doe",
	}, cont)

	assert.Equal(t, http.StatusOK, statusCode)

	mailer.WaitSingleMessageSent()
	verCode1 := parseVerificationCode(mailer.Msgs["jane_doe@gmail.com"][0])

	statusCode, resp2 := accounts.HandleAccountsRegister(&accounts.RegisterRequest{
		Login:     "jane_doe",
		Password:  "qwerty1",
		Email:     "jane_doe1@gmail.com",
		FirstName: "Jane1",
		LastName:  "Doe1",
	}, cont)

	assert.Equal(t, http.StatusOK, statusCode)

	mailer.WaitSingleMessageSent()
	verCode2 := parseVerificationCode(mailer.Msgs["jane_doe1@gmail.com"][0])

	statusCode, _ = accounts.HandleAccountsValidateEmail(&accounts.ValidateEmailRequest{
		VerificationCode: verCode2,
		Token:            resp2.Response.RandomToken,
	}, cont)

	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, _ = accounts.HandleAccountsValidateEmail(&accounts.ValidateEmailRequest{
		VerificationCode: verCode1,
		Token:            resp1.Response.RandomToken,
	}, cont)

	assert.Equal(t, http.StatusInternalServerError, statusCode)

	user := db.LoginToUser["jane_doe"]
	assert.Equal(t, "jane_doe", user.Login)
	assert.Equal(t, "qwerty1", user.Password)
	assert.Equal(t, "Jane1", user.FirstName)
	assert.Equal(t, "Doe1", user.LastName)
	assert.Equal(t, "jane_doe1@gmail.com", user.Email)
}

func TestChangePassword(t *testing.T) {
	cont, db, mailer := dependency.NewTestContext()

	ok, _ := registerUser(&models.User{
		Login:     "jane_doe",
		Password:  "qwerty",
		Email:     "jane_doe@gmail.com",
		FirstName: "Jane",
		LastName:  "Doe",
	}, cont, db, mailer)
	if !ok {
		t.Fatalf("Failed to create user")
	}

	statusCode, begResp := accounts.HandleAccountsBeginChangePassword(&accounts.BeginChangePasswordRequest{Email: "jane_doe@gmail.com"}, cont)
	assert.Equal(t, http.StatusOK, statusCode)
	idToken := begResp.Response.Token

	mailer.WaitSingleMessageSent()
	verCode := parseVerificationCode(mailer.GetLastEmail("jane_doe@gmail.com"))

	statusCode, commResp := accounts.HandleAccountsValidateChangePassword(&accounts.ValidateChangePasswordRequest{
		RandomToken:      idToken,
		VerificationCode: verCode,
	}, cont)

	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, _ = accounts.HandleAccountsCommitChangePassword(&accounts.CommitChangePasswordRequest{
		ChangePasswordToken: commResp.Response.ChangePasswordToken,
		NewPassword:         "asdfasdf",
	}, cont)

	assert.Equal(t, http.StatusOK, statusCode)

	passwd, _ := db.GetPassword("jane_doe")
	assert.Equal(t, "asdfasdf", passwd)
}

func TestChangePasswordWithRetry(t *testing.T) {
	cont, db, mailer := dependency.NewTestContext()

	ok, _ := registerUser(&models.User{
		Login:     "jane_doe",
		Password:  "qwerty",
		Email:     "jane_doe@gmail.com",
		FirstName: "Jane",
		LastName:  "Doe",
	}, cont, db, mailer)
	if !ok {
		t.Fatalf("Failed to create user")
	}

	statusCode, begResp := accounts.HandleAccountsBeginChangePassword(&accounts.BeginChangePasswordRequest{Email: "jane_doe@gmail.com"}, cont)
	assert.Equal(t, http.StatusOK, statusCode)
	idToken := begResp.Response.Token

	mailer.WaitSingleMessageSent()

	statusCode, begResp2 := accounts.HandleAccountsRevalidateChangePasswordCredentials(&accounts.RevalidateChangePasswordCredentialsRequest{RandomToken: idToken}, cont)
	assert.Equal(t, http.StatusOK, statusCode)

	idToken = begResp2.Response.RandomToken
	mailer.WaitSingleMessageSent()
	verCode := parseVerificationCode(mailer.GetLastEmail("jane_doe@gmail.com"))

	statusCode, commResp := accounts.HandleAccountsValidateChangePassword(&accounts.ValidateChangePasswordRequest{
		RandomToken:      idToken,
		VerificationCode: verCode,
	}, cont)

	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, _ = accounts.HandleAccountsCommitChangePassword(&accounts.CommitChangePasswordRequest{
		ChangePasswordToken: commResp.Response.ChangePasswordToken,
		NewPassword:         "asdfasdf",
	}, cont)

	assert.Equal(t, http.StatusOK, statusCode)

	passwd, _ := db.GetPassword("jane_doe")
	assert.Equal(t, "asdfasdf", passwd)
}

func registerUser(user *models.User, cont *dependency.DependencyContext, db *test_storage.TestDAO, mailer *dependency.TestMailer) (bool, uint32) {

	statusCode, regResp := accounts.HandleAccountsRegister(&accounts.RegisterRequest{
		Login:     user.Login,
		Password:  user.Password,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, cont)

	if statusCode != http.StatusOK {
		return false, 0
	}

	mailer.WaitSingleMessageSent()
	msg := mailer.Msgs["jane_doe@gmail.com"][0]
	verificationCode := parseVerificationCode(msg)
	statusCode, _ = accounts.HandleAccountsValidateEmail(&accounts.ValidateEmailRequest{
		VerificationCode: verificationCode,
		Token:            regResp.Response.RandomToken,
	}, cont)
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
