package pwd_hasher

import (
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hashedPassword), nil
}

func ComparePasswords(pwdHash string, pwd string) (bool, error) {
	bytes, err := hex.DecodeString(pwdHash)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(bytes, []byte(pwd))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
