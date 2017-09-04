package handlers

import (
	"bytes"
	"errors"
	entities "github.com/gost/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandler() (*entities.Thing, []error) {
	return nil, nil
}

func testHandlerError() (*entities.Thing, []error) {
	return nil, []error{errors.New("Test error")}
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
	assert.Equal(t, http.StatusBadRequest, rr.Code)
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
	assert.Equal(t, http.StatusBadRequest, rr.Code)
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
	assert.Equal(t, http.StatusBadRequest, rr.Code)
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
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandlePutTestWithGoodBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	thing.SetAllLinks("localhost")
	thing.ID = 1
	req, _ := http.NewRequest("PUT", "/bla", bytes.NewReader([]byte("{\"@iot.id\": 1, \"name\": \"thing1\", \"description\": \"test thing 1\"}")))
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, []error) { return testHandler() }

	// act
	handlePutRequest(rr, nil, req, thing, &handle, false)

	// assert

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, thing.GetSelfLink(), rr.HeaderMap.Get("Location"), "Expected header with Location to entity")
}
