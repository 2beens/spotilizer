package models

// ViewData object returned to frontend clients
type ViewData struct {
	Username string      `json:"username"`
	Title    string      `json:"title"`
	Message  string      `json:"message"`
	Error    string      `json:"error"`
	Data     interface{} `json:"data"`
}
