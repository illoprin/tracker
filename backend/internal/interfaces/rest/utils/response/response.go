package response

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "error"
)

func Error(message string) Response {
	return Response{
		Status:  StatusError,
		Message: message,
	}
}

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}
