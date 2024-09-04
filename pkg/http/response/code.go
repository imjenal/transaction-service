package response

type ErrorCode int

const (
	// DefaultErrorCode - in case other codes are irrelevant
	DefaultErrorCode ErrorCode = 500
	// EmptyRequestBody - when the request body is empty
	EmptyRequestBody ErrorCode = 1000
	// InvalidJSON - when the json data in request in invalid
	InvalidJSON ErrorCode = 1001
	// InvalidJSONField - when a json field in request in invalid
	InvalidJSONField ErrorCode = 1002
	// UnKnownJSONField - when json request body contains unnecessary field for the request
	UnKnownJSONField ErrorCode = 1003
	// RequestSizeExceeds - when request body's size is large
	RequestSizeExceeds ErrorCode = 1004
	// UnknownParseError - when the parse error is none of the above categories
	UnknownParseError ErrorCode = 1005
	// ValidationFailed - when the parse error is none of the above categories
	ValidationFailed ErrorCode = 1006
	//InvalidPathParam - when the path parameter is invalid
	InvalidPathParam ErrorCode = 1007
	//InvalidUUID - when the uuid is invalid
	InvalidUUID ErrorCode = 1008

	//ErrAccountNotFound - when account isn't found
	ErrAccountNotFound ErrorCode = 2001

	//ErrTransactionNotFound - when transaction isn't found
	ErrTransactionNotFound ErrorCode = 3001

	//ErrUserNotFound - when user isn't found
	ErrUserNotFound ErrorCode = 4001

	//ErrOperationTypeNotFound - when operation type isn't found
	ErrOperationTypeNotFound ErrorCode = 5001
)
