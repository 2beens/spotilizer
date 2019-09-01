package models

// ViewData object returned to frontend clients
type ViewData struct {
	Username string      `json:"username"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
}

type ErrorViewData struct {
	Username string      `json:"username"`
	Title    string      `json:"title"`
	Message  string      `json:"message"`
	Error    string      `json:"error"`
	Data     interface{} `json:"data"`
}

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
