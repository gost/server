package handlers

import (
	"bytes"
	entities "github.com/gost/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePostTestWithWrongDataType(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("POST", "/bla", nil)
	req.Header.Set("Content-Type", "this is an invalid content-type")
	handle := func() (interface{}, []error) { return testHandler() }

	// act
	handlePostRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandlePostTestWithNoBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("POST", "/bla", nil)
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, []error) { return testHandler() }

	// act
	handlePostRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandlePostTestWithWrongBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("POST", "/bla", bytes.NewReader([]byte("{\"name\": 10}")))
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, []error) { return testHandler() }

	// act
	handlePostRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandlePostTestWithError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("POST", "/bla", bytes.NewReader([]byte("{\"name\": \"thing1\", \"description\": \"test thing 1\"}")))
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, []error) { return testHandlerError() }

	// act
	handlePostRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandlePostTestWithGoodBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	thing.SetAllLinks("localhost")
	thing.ID = 1
	req, _ := http.NewRequest("POST", "/bla", bytes.NewReader([]byte("{ \"name\": \"thing1\", \"description\": \"test thing 1\"}")))
	req.Header.Set("Content-Type", "application/json")
	handle := func() (interface{}, []error) { return testHandler() }

	// act
	handlePostRequest(rr, nil, req, thing, &handle, false)

	// assert

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, thing.GetSelfLink(), rr.HeaderMap.Get("Location"), "Expected header with Location to entity")
}
