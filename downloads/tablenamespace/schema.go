package tablenamespace

type TableSearchQuery struct {
	From string `json:"from" form:"from"`
	To   string `json:"to" form:"to"`
}

type ResultResponse struct {
	Succeed bool   `json:"succeed"`
	Message string `json:"message,omitempty"`
}
