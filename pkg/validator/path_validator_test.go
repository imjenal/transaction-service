package validator

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/imjenal/transaction-service/pkg/http/response"
	"github.com/stretchr/testify/assert"
)

const (
	dummyUUID = "d99ee9ef-0507-4285-ba9e-a2b8ae2bb031"
)

func TestPathValidatorHandler(t *testing.T) {
	router := getTestRouter(t)

	t.Run("should return 200 for unknown keys", func(t *testing.T) {
		router.HandleFunc("/{userID}/{something}", fakeNextHandler)

		r := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/%s/%s", dummyUUID, "finvu"),
			nil,
		)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("should return 400 for invalid values", func(t *testing.T) {
		router.HandleFunc("/{userID}/{number}/{enum}", fakeNextHandler)

		r := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/%s/%d/%s", dummyUUID, 34, "no"),
			nil,
		)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		assert.Equal(t, 400, w.Code)
	})
}

func fakeNextHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("pass"))
	if err != nil {
		return
	}
}

func getTestRouter(t *testing.T) *mux.Router {
	t.Helper()
	jw := response.NewJSONWriter()
	v := New()
	pathMw := NewPathValidator(v, jw, map[string]string{
		"userID":  "uuid4",
		"fipUUID": "uuid4",
		"number":  "number",
		"enum":    "oneof=visa rupay mastercard",
	})
	router := mux.NewRouter()
	router.Use(pathMw)

	return router
}
