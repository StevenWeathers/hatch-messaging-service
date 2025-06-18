package http

import (
	"encoding/json"
	"net/http"
)

// sendJsonResponse adds the content type header for json and writes the response as encoded json
func sendJsonResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	jsonEncoder := json.NewEncoder(w)
	err := jsonEncoder.Encode(&response)

	return err
}
