package accounts

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/internal_data"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage/test_storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestBeginChangePassword(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		cont, db, mailer := dependency.NewTestContext()
		db.AddNewUser(0, "jane_doe", "qwerty", "jane.doe@gmail.com", "Jane", "Doe")
		code, resp := HandleAccountsBeginChangePassword(&BeginChangePasswordRequest{Email: "jane.doe@gmail.com"}, cont)

		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, 40, len(resp.Response.Token))

		data, ok := db.TokenToChangePassword[resp.Response.Token]
		if !ok {
			t.Fatalf("Failed to find change password data for token")
		}

		assert.Equal(t, "jane_doe", data.Login)
		assert.True(t, time.Now().Add(time.Minute*20).Before(*data.ExpireDate))
		assert.True(t, time.Now().Add(time.Minute*40).After(*data.ExpireDate))

		mailer.WaitSingleMessageSent()
		msg := mailer.GetLastEmail("jane.doe@gmail.com")
		assert.Contains(t, msg.Body, data.VerificationCode)
	})

	t.Run("Nonexistent email", func(t *testing.T) {
		cont, db, _ := dependency.NewTestContext()
		db.AddNewUser(0, "jane_doe", "qwerty", "jane.doe@gmail.com", "Jane", "Doe")
		code, resp := HandleAccountsBeginChangePassword(&BeginChangePasswordRequest{Email: "norma.jean@gmail.com"}, cont)

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, 0, len(db.TokenToChangePassword))
		assert.Equal(t, ErrorsUnion{
			IllegalEmail: &Error{},
		}, *resp.Errors)
	})
}

func TestCommitChangePassword(t *testing.T) {
	t.Run("Successful password change", func(t *testing.T) {
		cont, db, _ := dependency.NewTestContext()
		db.AddNewUser(0, "jane_doe", "qwerty", "jane.doe@gmail.com", "Jane", "Doe")
		expTime := time.Now().Add(5 * time.Second)
		db.AddChangePwdData(&test_storage.TestChangePasswordData{
			Token:                    "abcdefg",
			Login:                    "jane_doe",
			ExpireDate:               &expTime,
			VerificationCode:         "123456",
			ChangePasswordToken:      "qwerqwer",
			ChangePasswordExpireTime: &expTime,
		})

		code, _ := HandleAccountsCommitChangePassword(&CommitChangePasswordRequest{
			ChangePasswordToken: "qwerqwer",
			NewPassword:         "asdfasdf",
		}, cont)

		assert.Equal(t, http.StatusOK, code)
		user := db.LoginToUser["jane_doe"]

		assert.Equal(t, user.Password, "asdfasdf")
	})

	t.Run("Change password expired", func(t *testing.T) {
		cont, db, _ := dependency.NewTestContext()
		db.AddNewUser(0, "jane_doe", "qwerty", "jane.doe@gmail.com", "Jane", "Doe")
		expTime := time.Now().Add(-5 * time.Second)
		db.AddChangePwdData(&test_storage.TestChangePasswordData{
			Token:                    "abcdefg",
			Login:                    "jane_doe",
			ExpireDate:               &expTime,
			VerificationCode:         "123456",
			ChangePasswordToken:      "qwerqwer",
			ChangePasswordExpireTime: &expTime,
		})

		code, resp := HandleAccountsCommitChangePassword(&CommitChangePasswordRequest{
			ChangePasswordToken: "qwerqwer",
			NewPassword:         "asdfasdf",
		}, cont)

		assert.Equal(t, http.StatusBadRequest, code)
		user := db.LoginToUser["jane_doe"]

		assert.Equal(t, user.Password, "qwerty")
		assert.Equal(t, ErrorsUnion{
			ChangePasswordExpired: &Error{},
		}, *resp.Errors)
	})
}

func TestLogin(t *testing.T) {
	cont, db, _ := dependency.NewTestContext()
	db.AddNewUser(0, "jane_doe", "qwerty", "jane.doe@gmail.com", "Jane", "Doe")

	t.Run("Successful login", func(t *testing.T) {

		code, _ := HandleAccountsLogin(&LoginRequest{
			Login:    "jane_doe",
			Password: "qwerty",
		}, cont)

		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("Wrong password", func(t *testing.T) {

		code, resp := HandleAccountsLogin(&LoginRequest{
			Login:    "jane_doe",
			Password: "asdfasdfasdf",
		}, cont)

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, ErrorsUnion{
			InvalidLoginOrPassword: &Error{},
		}, *resp.Errors)
	})
}

func TestRegister(t *testing.T) {
	cont, db, _ := dependency.NewTestContext()

	t.Run("Successful registration", func(t *testing.T) {
		code, resp := HandleAccountsRegister(&RegisterRequest{
			Login:     "jane_doe",
			Password:  "qwerty",
			Email:     "jane.doe@gmail.com",
			FirstName: "Jane",
			LastName:  "Doe",
		}, cont)

		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, 40, len(resp.Response.RandomToken))

		regData, ok := db.TokenToRegistration[resp.Response.RandomToken]
		if !ok {
			t.Fatalf("User registration data not found")
		}
		assert.Equal(t, regData.Login, "jane_doe")
		assert.Equal(t, regData.Password, "qwerty")
		assert.Equal(t, regData.Email, "jane.doe@gmail.com")
		assert.Equal(t, regData.FirstName, "Jane")
		assert.Equal(t, regData.LastName, "Doe")
		assert.Equal(t, 6, len(regData.VerificationCode))
	})

	t.Run("Existing login error", func(t *testing.T) {
		db.AddNewUser(0, "jane_doe", "qwerty", "aaaaaaa@gmail.com", "Jane", "Doe")

		code, resp := HandleAccountsRegister(&RegisterRequest{
			Login:     "jane_doe",
			Password:  "qwerty",
			Email:     "jane.doe@gmail.com",
			FirstName: "Jane",
			LastName:  "Doe",
		}, cont)

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, ErrorsUnion{
			LoginIsNotAvailable: &Error{},
		}, *resp.Errors)
	})
}

func TestRevalidateChangePasswordCredentials(t *testing.T) {
	cont, db, mailer := dependency.NewTestContext()
	db.AddNewUser(0, "jane_doe", "qwerty", "aaaaaaa@gmail.com", "Jane", "Doe")

	expTime := time.Now().Add(internal_data.ResentCodeIn)
	db.AddChangePwdData(&test_storage.TestChangePasswordData{
		Token:                    "aaaaaa",
		Login:                    "jane_doe",
		ExpireDate:               &expTime,
		VerificationCode:         "123456",
		ChangePasswordToken:      "",
		ChangePasswordExpireTime: nil,
	})

	t.Run("Successful revalidation", func(t *testing.T) {

		code, resp := HandleAccountsRevalidateChangePasswordCredentials(
			&RevalidateChangePasswordCredentialsRequest{RandomToken: "aaaaaa"},
			cont)

		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, 40, len(resp.Response.RandomToken))

		mailer.WaitSingleMessageSent()
		assert.Equal(t, 1, len(mailer.Msgs))

		changeData, ok := db.TokenToChangePassword[resp.Response.RandomToken]
		if !ok {
			t.Fatalf("User change password data not found")
		}
		assert.Equal(t, changeData.Login, "jane_doe")
		assert.Equal(t, 6, len(changeData.VerificationCode))
		msg := mailer.Msgs["aaaaaaa@gmail.com"][0]
		assert.Contains(t, msg.Body, changeData.VerificationCode)
	})

	t.Run("Wrong token", func(t *testing.T) {

		code, resp := HandleAccountsRevalidateChangePasswordCredentials(
			&RevalidateChangePasswordCredentialsRequest{RandomToken: "aaaaab"},
			cont)

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, ErrorsUnion{
			InvalidChangePasswordRevalidation: &Error{},
		}, *resp.Errors)
	})
}

func TestRevalidateRegistrationCredentials(t *testing.T) {
	cont, db, mailer := dependency.NewTestContext()

	expTime := time.Now().Add(20 * time.Minute)
	db.AddRegDetails(&test_storage.TestRegistrationData{
		Token:            "aaaaaa",
		Login:            "jane_doe",
		Password:         "qwerty",
		Email:            "jane_doe@gmail.com",
		FirstName:        "Jane",
		LastName:         "Doe",
		ExpireData:       &expTime,
		VerificationCode: "123123",
	})

	t.Run("Successful revalidation", func(t *testing.T) {

		code, resp := HandleAccountsRevalidateRegistrationCredentials(
			&RevalidateRegistrationCredentialsRequest{Token: "aaaaaa"},
			cont)

		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, 40, len(resp.Response.RandomToken))

		mailer.WaitSingleMessageSent()
		assert.Equal(t, 1, len(mailer.Msgs))

		regData, ok := db.TokenToRegistration[resp.Response.RandomToken]
		if !ok {
			t.Fatalf("User change password data not found")
		}
		assert.Equal(t, regData.Login, "jane_doe")
		assert.Equal(t, 6, len(regData.VerificationCode))

		msg := mailer.Msgs["jane_doe@gmail.com"][0]
		assert.Contains(t, msg.Body, regData.VerificationCode)
	})

	t.Run("Wrong token", func(t *testing.T) {

		code, resp := HandleAccountsRevalidateRegistrationCredentials(
			&RevalidateRegistrationCredentialsRequest{Token: "aaaaab"},
			cont)

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, ErrorsUnion{
			InvalidRegistrationRevalidation: &Error{},
		}, *resp.Errors)
	})
}

func TestValidateChangePassword(t *testing.T) {
	cont, db, _ := dependency.NewTestContext()

	expTime := time.Now().Add(20 * time.Minute)
	db.AddChangePwdData(&test_storage.TestChangePasswordData{
		Token:            "aaaaaa",
		Login:            "jane_doe",
		ExpireDate:       &expTime,
		VerificationCode: "123123",
	})

	t.Run("Successful validation", func(t *testing.T) {

		code, resp := HandleAccountsValidateChangePassword(
			&ValidateChangePasswordRequest{RandomToken: "aaaaaa", VerificationCode: "123123"},
			cont)

		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, 40, len(resp.Response.ChangePasswordToken))

		regData, ok := db.TokenToChangePassword["aaaaaa"]
		if !ok {
			t.Fatalf("User change password data not found")
		}

		assert.Equal(t, 40, len(regData.ChangePasswordToken))

	})

	t.Run("Wrong verification code", func(t *testing.T) {

		code, resp := HandleAccountsValidateChangePassword(
			&ValidateChangePasswordRequest{RandomToken: "aaaaaa", VerificationCode: "111111"},
			cont)

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, ErrorsUnion{
			InvalidCode: &Error{},
		}, *resp.Errors)
	})
}

func TestValidateEmail(t *testing.T) {
	cont, db, _ := dependency.NewTestContext()

	expTime := time.Now().Add(20 * time.Minute)
	db.AddRegDetails(&test_storage.TestRegistrationData{
		Token:            "aaaaaa",
		Login:            "jane_doe",
		Password:         "qwerty",
		Email:            "jane_doe@gmail.com",
		FirstName:        "Jane",
		LastName:         "Doe",
		ExpireData:       &expTime,
		VerificationCode: "123123",
	})

	t.Run("Successful validation", func(t *testing.T) {

		code, _ := HandleAccountsValidateEmail(
			&ValidateEmailRequest{Token: "aaaaaa", VerificationCode: "123123"},
			cont)

		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("Wrong verification code", func(t *testing.T) {

		code, resp := HandleAccountsValidateEmail(
			&ValidateEmailRequest{Token: "aaaaaa", VerificationCode: "111111"},
			cont)

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, ErrorsUnion{
			InvalidCode: &Error{},
		}, *resp.Errors)
	})
}
