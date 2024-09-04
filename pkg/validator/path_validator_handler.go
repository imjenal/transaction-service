package validator

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imjenal/transaction-service/pkg/http/response"
)

// PathValidator is a middleware that validates the path variables
// It takes a map of variable names and tags to validate them
//
// For example, if you have a path like /users/{id} and you want to validate the id to be an uuid, you can do it like this:
//
//	variablesTagMap := map[string]string{
//		"id": "uuid",
//	}
type PathValidator struct {
	v               *Validator
	jsonWriter      *response.JSONWriter
	variablesTagMap map[string]string
	next            http.Handler
}

func NewPathValidator(v *Validator, jsonWriter *response.JSONWriter, variablesTagMap map[string]string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &PathValidator{
			v:               v,
			jsonWriter:      jsonWriter,
			variablesTagMap: variablesTagMap,
			next:            h,
		}
	}
}

func (p PathValidator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for variable, tags := range p.variablesTagMap {
		value := mux.Vars(r)[variable]
		if value == "" {
			continue
		}

		result, err := p.v.IsValidString(r.Context(), value, tags)
		if err != nil {
			fmt.Println("PathValidator.ServeHTTP: failed to validate string")
			p.jsonWriter.DefaultError(w)

			return
		}

		if !result.Valid {
			data := map[string]string{
				"variable": variable,
				"value":    value,
			}
			p.jsonWriter.BadRequest(w, response.NewError(response.InvalidPathParam, "Invalid path param", "send correct path param", data))

			return
		}
	}

	p.next.ServeHTTP(w, r)
}
