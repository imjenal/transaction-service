package response

const (
	defaultErrMsg = "Unknown error occurred"
	defaultFix    = "Please try again in sometime. If the issue persists, please contact us."
)

type (
	// APIError is to be implemented for sending an error to the client
	APIError struct {
		//Code is 3-4 digit error code that informs client about the error
		Code ErrorCode `json:"code"`

		//Message is the error that happened.
		Message string `json:"message"`

		//HowToFix is to be shown to the user on how to fix the error. This field is not sent when empty
		HowToFix string `json:"how_to_fix,omitempty"`

		//Data can be added when there's any additional data on why error occurred,
		// or can contain the field names that failed validation
		Data interface{} `json:"data"`
	}

	response struct {
		Data  interface{} `json:"data"`
		Error *APIError   `json:"error"`
	}
)

// NewError creates a new JSON error response
func NewError(
	code ErrorCode,
	message string,
	howToFix string,
	data interface{},
) *APIError {
	return &APIError{
		Code:     code,
		Message:  message,
		HowToFix: howToFix,
		Data:     data,
	}
}

var DefaultErr = NewError(DefaultErrorCode, defaultErrMsg, defaultFix, nil)
