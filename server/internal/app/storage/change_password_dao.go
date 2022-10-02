package storage

import "time"

type ChangePasswordDao struct {
	Storage *Storage
}

// GetVerificationCodeInfo returns (if verification code is valid, expiresAt, error). Returns empty string if not found
func (cpd *ChangePasswordDao) GetVerificationCodeInfo(identificationToken string) (string, *time.Time, error) {
	row := cpd.Storage.Db.QueryRow(`
        SELECT (verification_code, expire_date)
        FROM change_password_base
        WHERE token = $1
    `, identificationToken)

	var expiresAt = time.Time{}
	var verificationCode string
	err := row.Scan(&verificationCode, &expiresAt)
	if err != nil {
		return "", &expiresAt, err
	}
	return verificationCode, &expiresAt, nil
}

// SetChangePasswordToken remembers changePasswordToken for identificationToken issued in accounts/begin_change_password
func (cpd *ChangePasswordDao) SetChangePasswordToken(identificationToken, changePasswordToken string) error {
	_, err := cpd.Storage.Db.Exec(`
	        UPDATE change_password_base 
	        SET change_password_token = $1
	        WHERE token = $2
	    `, changePasswordToken, identificationToken)
	return err
}