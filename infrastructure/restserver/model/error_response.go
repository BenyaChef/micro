package restservermodel

type ErrorResponse struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{Code: code, Error: message}
}
