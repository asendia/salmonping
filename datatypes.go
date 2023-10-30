package main

import "net/http"

type DefaultResponse struct {
	Message string `json:"message"`
}

type DefaultErrorResponse struct {
	Error   string      `json:"error"`
	Header  http.Header `json:"header"`
	Level   string      `json:"level"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
	Query   string      `json:"query"`
}

func (resp *DefaultErrorResponse) JSON() map[string]interface{} {
	return map[string]interface{}{
		"error":   resp.Error,
		"header":  resp.Header,
		"level":   resp.Level,
		"message": resp.Message,
		"payload": resp.Payload,
		"query":   resp.Query,
	}
}
