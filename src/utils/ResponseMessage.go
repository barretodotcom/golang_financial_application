package utils

import "encoding/json"

type ResponseMessage struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func NewResponseMessage(message string, status string) []byte {

	bytesResponse, _ := json.Marshal(ResponseMessage{
		Message: message,
		Status:  status,
	})

	return bytesResponse
}
