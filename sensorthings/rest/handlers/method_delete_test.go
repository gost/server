package handlers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandlerDelete() error {
	return nil
}

func testHandlerDeleteError() error {
	return errors.New("Test error")
}

func TestHandleDeleteWithError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/bla", nil)
	handle := func() error { return testHandlerDeleteError() }

	// act
	handleDeleteRequest(rr, nil, req, &handle, false)

	// assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandleDeleteTestOk(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/bla", nil)
	handle := func() error { return testHandlerDelete() }

	// act
	handleDeleteRequest(rr, nil, req, &handle, false)

	// assert

	assert.Equal(t, http.StatusOK, rr.Code)
}
