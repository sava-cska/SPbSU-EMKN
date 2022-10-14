package internal_data

func ValidatePassword(password string) bool {
	return len(password) != 0
}
