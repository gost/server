package errors

import (
	"net/http"
)

// APIError holds information about an error including status codes.
type APIError struct {
	error          error
	httpStatusCode int
}

// GetHTTPErrorStatusCode returns the status code.
func (e APIError) GetHTTPErrorStatusCode() int {
	return e.httpStatusCode
}

// Error implements the error interface for apiError
func (e APIError) Error() string {
	return e.error.Error()
}

// NewErrorWithStatusCode creates a new apiError with a given status code
func NewErrorWithStatusCode(err error, status int) error {
	return APIError{err, status}
}

// NewBadRequestError creates an apiError with status code 400.
func NewBadRequestError(err error) error {
	return NewErrorWithStatusCode(err, http.StatusBadRequest)
}

// NewConflictRequestError creates an apiError with status code 409.
func NewConflictRequestError(err error) error {
	return NewErrorWithStatusCode(err, http.StatusConflict)
}

// NewRequestNotImplemented creates an apiError with status code 501.
func NewRequestNotImplemented(err error) error {
	return NewErrorWithStatusCode(err, http.StatusNotImplemented)
}

// NewRequestNotFound creates an apiError with status code 404.
func NewRequestNotFound(err error) error {
	return NewErrorWithStatusCode(err, http.StatusNotFound)
}

// NewRequestMethodNotAllowed creates an apiError with status code 405.
func NewRequestMethodNotAllowed(err error) error {
	return NewErrorWithStatusCode(err, http.StatusMethodNotAllowed)
}

// NewRequestInternalServerError creates an apiError with status code 500.
func NewRequestInternalServerError(err error) error {
	return NewErrorWithStatusCode(err, http.StatusInternalServerError)
}
