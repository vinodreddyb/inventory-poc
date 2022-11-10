package responses

type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Body    any    `json:"body"`
}
