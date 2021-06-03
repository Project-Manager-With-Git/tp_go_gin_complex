package usernamespace

type UserCreateQuery struct {
	Name string `json:"Name"`
}

type UserUpdateQuery struct {
	Name string `json:"Name"`
}

type ResultResponse struct {
	Succeed bool   `json:"succeed"`
	Message string `json:"message,omitempty"`
}
type LinkResponse struct {
	URI         string `json:"uri"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

type UserListResponse struct {
	Description string         `json:"Description"`
	UserCount   int64          `json:"UserCount"`
	Links       []LinkResponse `json:"Links"`
}
