package httphelper

type APIResponse struct {
	Data interface{} `json:"data,omitempty"`
	Meta interface{} `json:"meta,omitempty"`
	Err  interface{} `json:"error,omitempty"`
}

func Success(data interface{}) APIResponse {
	return APIResponse{Data: data}
}

func SuccessWithMeta(data, meta interface{}) APIResponse {
	return APIResponse{Data: data, Meta: meta}
}

func Error(msg string) APIResponse {
	return APIResponse{Err: map[string]string{"message": msg}}
}
