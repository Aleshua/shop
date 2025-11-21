package responses

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func NewResponse(msg string, data, meta interface{}) Response {
	return Response{Message: msg, Data: data, Meta: meta}
}
