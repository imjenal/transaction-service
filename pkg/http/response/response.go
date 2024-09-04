package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
)

type JSONWriter struct{}

// NewJSONWriter creates a new instance of JSONWriter
func NewJSONWriter() *JSONWriter {
	return &JSONWriter{}
}

// Ok sends the data to client with http status 200
func (j *JSONWriter) Ok(w http.ResponseWriter, data interface{}) {
	res := j.buildResponse(data, nil)
	j.jsonWrite(w, res, http.StatusOK)
}

// Error sends error to client with the given http status
func (j *JSONWriter) Error(w http.ResponseWriter, apiError *APIError, httpStatus int) {
	res := j.buildResponse(nil, apiError)
	j.jsonWrite(w, res, httpStatus)
}

// NotFound sends error to client with http status 404
func (j *JSONWriter) NotFound(w http.ResponseWriter, apiError *APIError) {
	j.Error(w, apiError, http.StatusNotFound)
}

// Unauthorized sends error to client with http status 401
func (j *JSONWriter) Unauthorized(w http.ResponseWriter, apiError *APIError) {
	j.Error(w, apiError, http.StatusUnauthorized)
}

// Forbidden sends error to client with http status 403
func (j *JSONWriter) Forbidden(w http.ResponseWriter, r *http.Request, apiError *APIError) {
	j.Error(w, apiError, http.StatusForbidden)
}

// UnprocessableEntity sends error to client with http status 422
func (j *JSONWriter) UnprocessableEntity(w http.ResponseWriter, apiError *APIError) {
	j.Error(w, apiError, http.StatusUnprocessableEntity)
}

// BadRequest sends error to client with http status 400
func (j *JSONWriter) BadRequest(w http.ResponseWriter, apiError *APIError) {
	j.Error(w, apiError, http.StatusBadRequest)
}

// Internal sends error to client with http status 500
func (j *JSONWriter) Internal(w http.ResponseWriter, apiError *APIError) {
	j.Error(w, apiError, http.StatusInternalServerError)
}

// TooManyRequest sends error to client with http status 429
func (j *JSONWriter) TooManyRequest(w http.ResponseWriter, apiError *APIError) {
	j.Error(w, apiError, http.StatusTooManyRequests)
}

func (j *JSONWriter) Conflict(w http.ResponseWriter, ae *APIError) {
	j.Error(w, ae, http.StatusConflict)
}

// DefaultError sends an unknown error to theclient with HTTP status 500
func (j *JSONWriter) DefaultError(w http.ResponseWriter) {
	j.Error(w, DefaultErr, http.StatusInternalServerError)
}

// buildResponse builds the response object to be sent to the client
func (j *JSONWriter) buildResponse(data interface{}, apiError *APIError) *response {
	return &response{
		Data:  data,
		Error: apiError,
	}
}

// jsonWrite writes the data to the client as JSON
func (j *JSONWriter) jsonWrite(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)

	if err != nil && errors.Is(err, os.ErrDeadlineExceeded) {
		log.Printf("jsonWrite: failed to write response due to io timeout: %v", err)
		return
	}

	if err != nil {
		log.Printf("jsonWrite: failed to write response: %v\n", err)
	}
}
