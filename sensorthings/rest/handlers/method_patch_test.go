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

func testHandlerPatch() (*entities.Thing, error) {
	return nil, nil
}

func testHandlerPatchError() (*entities.Thing, error) {
	return nil, errors.New("Test error")
}

func TestHandlePatchTestWithWrongDataType(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("PATCH", "/bla", nil)
	req.Header.Set("Content-Type", "this is an invalid content-type")
	handle := func() (interface{}, error) { return testHandlerPatch() }

	// act
	handlePatchRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t,  http.StatusBadRequest, rr.Code)
}

func TestHandlePatchTestWithNoBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("PATCH", "/bla", nil)
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, error) { return testHandlerPatch() }

	// act
	handlePatchRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t,  http.StatusBadRequest, rr.Code)
}

func TestHandlePatchTestWithWrongBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("PATCH", "/bla", bytes.NewReader([]byte("{\"name\": 10}")))
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, error) { return testHandlerPatch() }

	// act
	handlePatchRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t,  http.StatusBadRequest, rr.Code)
}

func TestHandlePatchTestWithPutError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("PATCH", "/bla", bytes.NewReader([]byte("{\"name\": \"thing1\", \"description\": \"test thing 1\"}")))
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, error) { return testHandlerPatchError() }

	// act
	handlePatchRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t,  http.StatusInternalServerError, rr.Code)
}

func TestHandlePatchTestWithGoodBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	thing.SetAllLinks("localhost")
	thing.ID = 1
	req, _ := http.NewRequest("PATCH", "/bla", bytes.NewReader([]byte("{\"@iot.id\": 1, \"name\": \"thing1\", \"description\": \"test thing 1\"}")))
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, error) { return testHandlerPatch() }

	// act
	handlePatchRequest(rr, nil, req, thing, &handle, false)

	// assert

	assert.Equal(t,  http.StatusOK, rr.Code)
	assert.Equal(t,  thing.GetSelfLink(), rr.HeaderMap.Get("Location"), "Expected header with Location to entity")
}
