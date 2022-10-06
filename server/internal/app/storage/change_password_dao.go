package storage

import (
	"database/sql"
	"time"
)

type ChangePasswordDao struct {
	Storage *Storage
}

// GetVerificationCodeInfo returns (if verification code is valid, expiresAt, error). Returns empty string if not found
func (cpd *ChangePasswordDao) GetVerificationCodeInfo(identificationToken string) (string, *time.Time, error) {
	row := cpd.Storage.Db.QueryRow(
		`SELECT verification_code, expire_date
               FROM change_password_base
               WHERE token = $1`, identificationToken)

	var expiresAt = time.Time{}
	var verificationCode string
	if err := row.Scan(&verificationCode, &expiresAt); err != nil {
		return "", &expiresAt, err
	}
	return verificationCode, &expiresAt, nil
}

// SetChangePasswordToken remembers changePasswordToken for identificationToken issued in accounts/begin_change_password
func (cpd *ChangePasswordDao) SetChangePasswordToken(identificationToken string, changeTime time.Time,
	changePasswordToken string) error {
	_, err := cpd.Storage.Db.Exec(
		`UPDATE change_password_base 
	           SET
			   change_password_token = $1,
			   change_password_expire_time = $2
	           WHERE token = $3`,
		changePasswordToken,
		changeTime,
		identificationToken,
	)
	return err
}

func (cpd *ChangePasswordDao) Upsert(token string, login string, expiredTime time.Time,
	verificationCode string) error {
	_, err := cpd.Storage.Db.Exec(
		`
		INSERT INTO change_password_base
		(token, login, expire_date, verification_code, change_password_token, change_password_expire_time)
		VALUES ($1, $2, $3, $4, $5, $6)
		`,
		token, login, expiredTime, verificationCode, token, time.Time{})
	return err
}

func (cpd *ChangePasswordDao) UpdateVerificationCode(token string, newVerificationCode string) (string, bool, error) {
	res := cpd.Storage.Db.QueryRow(
		`UPDATE change_password_base
               SET verification_code = $1
               WHERE token = $2
               RETURNING login`,
		newVerificationCode, token,
	)
	login := ""
	if err := res.Scan(&login); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		} else {
			return "", false, err
		}
	}
	return login, true, nil
}

func (cpd *ChangePasswordDao) FindPwdToken(changePwdToken string) (string, time.Time, error) {
	changePwdRecord := cpd.Storage.Db.QueryRow(
		`
		SELECT login, change_password_expire_time
		FROM change_password_base
		WHERE change_password_token = $1
		`,
		changePwdToken)

	var login string
	var expiredTime time.Time
	if errScan := changePwdRecord.Scan(&login, &expiredTime); errScan != nil {
		return "", time.Time{}, errScan
	}
	return login, expiredTime, nil
}
