package storage

import (
	"time"
)

type ChangePasswordDAO struct {
	Storage *Storage
}

// GetVerificationCodeInfo returns (if verification code is valid, expiresAt, error). Returns empty string if not found
func (cpd *ChangePasswordDAO) GetVerificationCodeInfo(identificationToken string) (string, *time.Time, error) {
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
func (cpd *ChangePasswordDAO) SetChangePasswordToken(identificationToken string, changeTime time.Time,
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

func (cpd *ChangePasswordDAO) Upsert(token string, login string, expiredTime time.Time,
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

func (cpd *ChangePasswordDAO) FindTokenAndDelete(token string) (string, error) {
	tx, err := cpd.Storage.Db.Begin()
	if err != nil {
		return "", err
	}
	registerRecord := tx.QueryRow(`SELECT login FROM change_password_base WHERE token = $1`, token)

	var login string
	if errScan := registerRecord.Scan(&login); errScan != nil {
		tx.Rollback()
		return "", errScan
	}

	_, err = tx.Exec(`DELETE FROM change_password_base WHERE token = $1`, token)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return login, nil
}

func (cpd *ChangePasswordDAO) FindPwdToken(changePwdToken string) (string, time.Time, error) {
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
