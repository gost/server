package handlers

import (
	"github.com/gost/server/sensorthings/entities"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"errors"
)

func testHandler() (*entities.Thing, []error) {
	return nil, nil
}

func testHandlerError() (*entities.Thing, []error) {
	return nil, []error {errors.New("Test error")}
}

func TestHandlePutTestWithWrongDataType(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("PUT", "/bla", nil)
	req.Header.Set("Content-Type", "this is an invalid content-type")
	handle := func() (interface{}, []error) { return testHandler() }

	// act
	handlePutRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t,  http.StatusBadRequest, rr.Code)
}

func TestHandlePutTestWithNoBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("PUT", "/bla", nil)
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, []error) { return testHandler() }

	// act
	handlePutRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t,  http.StatusBadRequest, rr.Code)
}

func TestHandlePutTestWithWrongBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("PUT", "/bla", bytes.NewReader([]byte("{\"name\": 10}")))
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, []error) { return testHandler() }

	// act
	handlePutRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t,  http.StatusBadRequest, rr.Code)
}

func TestHandlePutTestWithPutError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("PUT", "/bla", bytes.NewReader([]byte("{\"name\": \"thing1\", \"description\": \"test thing 1\"}")))
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, []error) { return testHandlerError() }

	// act
	handlePutRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t,  http.StatusInternalServerError, rr.Code)
}

func TestHandlePutTestWithGoodBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("PUT", "/bla", bytes.NewReader([]byte("{\"name\": \"thing1\", \"description\": \"test thing 1\"}")))
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, []error) { return testHandler() }

	// act
	handlePutRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t,  http.StatusOK, rr.Code)
}
