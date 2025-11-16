package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "ok"
	StatusError = "error"
)

func Ok() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(errMsg string) Response {
	return Response{
		Status: StatusError,
		Error:  errMsg,
	}
}

// TODO: add validation error response with help message
