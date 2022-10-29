package profiles

type Profile struct {
	Id        uint32 `json:"id"`
	AvatarUrl string `json:"avatar_url"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type GetWrapper struct {
	Profiles *[]Profile `json:"profiles,omitempty"`
}

type GetResponse struct {
	Errors   *ErrorsUnion `json:"errors,omitempty"`
	Response *GetWrapper  `json:"response,omitempty"`
}

func (response GetResponse) Bind() {}
