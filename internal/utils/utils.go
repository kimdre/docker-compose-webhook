package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonResponse struct {
	Code    int    `json:"code"`
	JobID   string `json:"job_id,omitempty"`
	Details string `json:"details,omitempty"`
}

// jsonError inherits from jsonResponse and adds an error message
type jsonError struct {
	jsonResponse
	Error string `json:"error"`
}

// JSONError writes an error response to the client in JSON format
func JSONError(w http.ResponseWriter, err interface{}, details, jobId string, code int) {
	if _, ok := err.(error); ok {
		err = fmt.Sprintf("%v", err)
	}

	resp := jsonError{
		Error: err.(string),
		jsonResponse: jsonResponse{
			Code:    code,
			JobID:   jobId,
			Details: details,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func JSONResponse(w http.ResponseWriter, details, jobId string, code int) {
	resp := jsonResponse{
		Code:    code,
		JobID:   jobId,
		Details: details,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}
