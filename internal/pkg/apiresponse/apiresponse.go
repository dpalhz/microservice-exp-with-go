package apiresponse

// Response defines a consistent API response structure.
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) Response {
	return Response{Status: "success", Data: data}
}

func Error(msg string) Response {
	return Response{Status: "error", Message: msg}
}
