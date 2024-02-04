package controllers

type ErrorResponse struct {
	Message string `json:"message"`
}

func Error(err error) *ErrorResponse {
	return &ErrorResponse{Message: err.Error()}
}
