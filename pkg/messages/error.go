package messages

type ErrorsResponse struct {
	Errors []ErrorMessage `json:"errors"`
}

type ErrorMessage struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}
