package actions

const (
	IllegalEmail        string = "ILLEGAL_EMAIL"
	LoginIsNotAvailable        = "LOGIN_IS_NOT_AVAILABLE"
	IllegalPassword            = "ILLEGAL_PASSWORD"
)

type Error struct {
	Code string `json:"code"`
}

type EmknResponse struct {
	Errors []Error `json:"errors"`
}
