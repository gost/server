package handlers

import (
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandler() (*entities.Thing, []error) {
	return nil, nil
}

func TestHandlePutTestWithWrongDataType(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("GET", "/bla", nil)
	req.Header.Set("Content-Type", "this is an invalid content-type")
	handle := func() (interface{}, []error) { return testHandler() }

	// act
	handlePutRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.True(t, rr.Code == http.StatusBadRequest)
}

func TestHandlePutTestWithNoBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	req, _ := http.NewRequest("GET", "/bla", nil)
	handle := func() (interface{}, []error) { return testHandler() }

	// act
	handlePutRequest(rr, nil, req, thing, &handle, false)

	// assert
	assert.True(t, rr.Code == http.StatusBadRequest)
}
