package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
	"github.com/imjenal/transaction-service/pkg/http/response"
	"github.com/imjenal/transaction-service/pkg/validator"
)

// Reader has functions to read, parse nad validate request data
type Reader struct {
	jw            *response.JSONWriter
	validator     *validator.Validator
	schemaDecoder *schema.Decoder
}

// NewReader returns a new instance of Reader
func NewReader(jw *response.JSONWriter, validator *validator.Validator) *Reader {
	return &Reader{
		jw:            jw,
		validator:     validator,
		schemaDecoder: schema.NewDecoder(),
	}
}

// ReadJSONAndValidate reads a json request body into the given struct and the validates the struct data
func (read *Reader) ReadJSONAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := read.ReadJSONRequest(r, v); err != nil {
		parseErr := read.HandleParseError(err)
		read.jw.BadRequest(w, parseErr.APIError())

		return false
	}

	ve := read.validate(r.Context(), v)
	if ve != nil {
		read.jw.UnprocessableEntity(w, ve)
		return false
	}

	return true
}

// ReadQueryParamsAndValidate reads the given query params to the given struct and the validates the struct data
func (read *Reader) ReadQueryParamsAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	err := read.ReadQueryParams(r, v)
	if err != nil {
		parseErr := read.HandleParseError(err)
		read.jw.BadRequest(w, parseErr.APIError())

		return false
	}

	ve := read.validate(r.Context(), v)
	if ve != nil {
		read.jw.UnprocessableEntity(w, ve)
		return false
	}

	return true
}

// ReadJSONRequest reads a json request body into the given struct
func (read *Reader) ReadJSONRequest(r *http.Request, v interface{}) error {
	var buf bytes.Buffer
	bodyCopy := io.TeeReader(r.Body, &buf)

	// We need to read the teed reader first
	// or else buf will be empty
	d := json.NewDecoder(bodyCopy)
	d.DisallowUnknownFields()

	if err := d.Decode(v); err != nil {
		return fmt.Errorf("error decoding request body: %v", err)
	}

	bts, err := io.ReadAll(&buf)
	if err != nil {
		return fmt.Errorf("error reading request body: %v", err)
	}

	rawJSON := removeAllWhitespace(string(bts))

	// Handle empty JSON as empty request body
	if rawJSON == "{}" {
		return io.EOF
	}

	return nil
}

// ReadQueryParams reads the given query params to the given struct
func (read *Reader) ReadQueryParams(r *http.Request, v interface{}) error {
	return read.schemaDecoder.Decode(v, r.URL.Query())
}

// validate functions uses the validator to test issues with the given data
func (read *Reader) validate(ctx context.Context, v interface{}) *response.APIError {
	result, err := read.validator.IsValidStruct(ctx, v)
	if err != nil {
		panic(err)
	}

	if result.Valid && len(result.Fields) == 0 {
		return nil
	}

	return response.NewError(
		response.ValidationFailed,
		"Invalid data received for request",
		"Please send data in correct format",
		result.Fields,
	)
}

func removeAllWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\t' || r == '\r' {
			return -1
		}

		return r
	}, str)
}
