package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/imjenal/transaction-service/pkg/http/response"
)

type (
	// ParseError struct defines the structure of the JSON parsing error
	ParseError struct {
		Msg     string             `json:"message"`
		ErrCode response.ErrorCode `json:"code"`
		Fix     string             `json:"fix"`
		ErrData interface{}        `json:"data"`
	}
)

func (p *ParseError) APIError() *response.APIError {
	return response.NewError(p.ErrCode, p.Msg, p.Fix, p.ErrData)
}

// HandleParseError checks the type of the error in request parsing and writes an appropriate response
func (read *Reader) HandleParseError(err error) *ParseError {
	var (
		syntaxError        *json.SyntaxError
		unmarshalTypeError *json.UnmarshalTypeError
	)

	switch {
	// Catch any syntax errors in the JSON and send an error Message
	// which interpolates the location of the problem to make it easier for the client to fix.
	// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error for syntax errors in the JSON.
	case errors.As(err, &syntaxError), errors.Is(err, io.ErrUnexpectedEOF):
		return &ParseError{
			Msg:     "Request body contains badly-formed JSON",
			Fix:     "Make sure the syntax of the JSON payload passed in the request body is valid.",
			ErrCode: response.InvalidJSON,
			ErrData: map[string]interface{}{
				"offset": syntaxError.Offset,
			},
		}

	// Catch any type errors, like trying to assign a string in the JSON request body to an int field in our Person struct.
	// We can interpolate the relevant field name and position into the error
	// Message to make it easier for the client to fix.
	case errors.As(err, &unmarshalTypeError):
		return &ParseError{
			Msg:     fmt.Sprintf("Request body contains an invalid value for the %q field", unmarshalTypeError.Field),
			Fix:     fmt.Sprintf("Make sure that the value passed to the field %q is a valid %q", unmarshalTypeError.Field, unmarshalTypeError.Type.Name()),
			ErrCode: response.InvalidJSONField,
			ErrData: map[string]interface{}{
				"field": unmarshalTypeError.Field,
			},
		}

	// Catch the error caused by extra unexpected fields in the request body.
	// We extract the field name from the error Message and interpolate it in our custom error Message.
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")

		return &ParseError{
			Msg:     fmt.Sprintf("Request body contains unknown field %s", fieldName),
			Fix:     fmt.Sprintf("Please remove the unknown field %s from the request body", fieldName),
			ErrCode: response.UnKnownJSONField,
			ErrData: map[string]interface{}{
				"field": fieldName,
			},
		}

	// An io.EOF error is returned by Decode() if the request body is empty.
	case errors.Is(err, io.EOF):
		return &ParseError{
			Msg:     "Request body must not be empty",
			Fix:     "Please make sure that the request body is not empty.",
			ErrCode: response.EmptyRequestBody,
			ErrData: nil,
		}

	// Catch the error caused by the request body being too large.
	case err.Error() == "http: request body too large":
		return &ParseError{
			Msg:     "Request body must not be larger than 1MB",
			Fix:     "Please make sure that the size of the request body payload is less than 1MB",
			ErrCode: response.RequestSizeExceeds,
			ErrData: nil,
		}
	}

	return &ParseError{
		Msg:     "Unknown error occurred",
		ErrCode: response.UnknownParseError,
		ErrData: nil,
	}
}
