package web

type WebResponseSuccess struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type WebResponseFailed struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
