package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testErrorMsg  = "error occurred"
	testErrorCode = DefaultErrorCode
	testErrorData = 1
)

var testError = NewError(testErrorCode, testErrorMsg, "", testErrorData)

func TestJSONWriter_Ok(t *testing.T) {
	jw := NewJSONWriter()
	rr, _ := getResponseRequest()

	jw.Ok(rr, nil)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestJSONWriter_BadRequest(t *testing.T) {
	jw := NewJSONWriter()
	rr, _ := getResponseRequest()

	jw.BadRequest(rr, testError)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestJSONWriter_DefaultError(t *testing.T) {
	jw := NewJSONWriter()
	rr, _ := getResponseRequest()

	jw.DefaultError(rr)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestJSONWriter_Forbidden(t *testing.T) {
	jw := NewJSONWriter()
	rr, r := getResponseRequest()

	jw.Forbidden(rr, r, testError)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestJSONWriter_Internal(t *testing.T) {
	jw := NewJSONWriter()
	rr, _ := getResponseRequest()

	jw.Internal(rr, testError)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestJSONWriter_NotFound(t *testing.T) {
	jw := NewJSONWriter()
	rr, _ := getResponseRequest()

	jw.NotFound(rr, testError)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestJSONWriter_Unauthorized(t *testing.T) {
	jw := NewJSONWriter()
	rr, _ := getResponseRequest()

	jw.Unauthorized(rr, testError)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJSONWriter_UnprocessableEntity(t *testing.T) {
	jw := NewJSONWriter()
	rr, _ := getResponseRequest()

	jw.UnprocessableEntity(rr, testError)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestJSONWriter_buildResponse(t *testing.T) {
	jw := NewJSONWriter()

	t.Run("test success response", func(t *testing.T) {
		res := jw.buildResponse(1, nil)
		assert.Nil(t, res.Error)
		assert.Equal(t, 1, res.Data)
	})

	t.Run("test failure response", func(t *testing.T) {
		res := jw.buildResponse(nil, testError)
		assert.Nil(t, res.Data)
	})
}

func TestJSONWriter_jsonWrite(t *testing.T) {

	jw := NewJSONWriter()
	rr, _ := getResponseRequest()

	data := map[string]string{
		"hello": "world",
	}

	jw.jsonWrite(rr, data, http.StatusOK)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func getResponseRequest() (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	return w, getRequest()
}

func getRequest() *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	return r
}
