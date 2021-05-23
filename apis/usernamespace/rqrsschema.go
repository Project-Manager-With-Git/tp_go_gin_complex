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
