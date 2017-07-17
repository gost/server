package errors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiError(t *testing.T) {
	// arrange
	err := errors.New("yo")
	var apierr = APIError{err, 200}

	// assert
	assert.Equal(t, 200, apierr.GetHTTPErrorStatusCode(), "should return 200")
	assert.Equal(t, "yo", apierr.Error(), "should return 200")
}

func TestRequestStatusCodes(t *testing.T) {
	// arrange
	badrequesterror := NewBadRequestError(errors.New("bad"))
	notfounderror := NewRequestNotFound(errors.New("notfound"))
	conflicterror := NewConflictRequestError(errors.New("conflict"))
	notimplementederror := NewRequestNotImplemented(errors.New("notimplemented"))
	notallowederror := NewRequestMethodNotAllowed(errors.New("notallowed"))
	internalservererror := NewRequestInternalServerError(errors.New("internalserver"))

	// assert
	assert.Equal(t, "bad", badrequesterror.Error())
	assert.Equal(t, "notfound", notfounderror.Error())
	assert.Equal(t, "conflict", conflicterror.Error())
	assert.Equal(t, "notimplemented", notimplementederror.Error())
	assert.Equal(t, "notallowed", notallowederror.Error())
	assert.Equal(t, "internalserver", internalservererror.Error())

}
