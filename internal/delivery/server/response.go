package server

import (
	"encoding/json"
	"net/http"

	"finance-management/tools/helpers"
)

func writeSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(helpers.NewApiSuccessResponse(data).ParseToByte())
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(helpers.NewApiErrorResponse(message, details).ParseToByte())
}

func decodeJSON(r *http.Request, destination interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(destination)
}
