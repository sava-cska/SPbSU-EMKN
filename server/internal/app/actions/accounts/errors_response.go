package accounts

type Error struct{}

type ErrorsUnion struct {
	IllegalPassword        *Error `json:"illegal_password,omitempty"`
	IllegalLogin           *Error `json:"illegal_login,omitempty"`
	IllegalEmail           *Error `json:"illegal_email,omitempty"`
	LoginIsNotAvailable    *Error `json:"login_is_not_available,omitempty"`
	RegistrationExpired    *Error `json:"registration_expired,omitempty"`
	InvalidCode            *Error `json:"code_invalid,omitempty"`
	ChangePasswordExpired  *Error `json:"password_change_expired,omitempty"`
	InvalidLoginOrPassword *Error `json:"invalid_login_or_password,omitempty"`
}
