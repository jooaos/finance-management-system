package helpers

import "encoding/json"

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

func NewApiSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Data: data,
	}
}

func (s SuccessResponse) ParseToByte() []byte {
	data, err := json.Marshal(s)
	if err != nil {
		return []byte{}
	}

	return data
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func NewApiErrorResponse(message string, details interface{}) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Details: details,
	}
}

func (e ErrorResponse) ParseToByte() []byte {
	data, err := json.Marshal(e)
	if err != nil {
		return []byte{}
	}

	return data
}
