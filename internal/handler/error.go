package handler

// APIError
type APIError struct {
	Message    string
	StatusCode int
}

func NewAPIError(message string, status int) APIError {
	return APIError{
		Message:    message,
		StatusCode: status,
	}
}

// Error implements error interface
func (a APIError) Error() string {
	return a.Message
}
