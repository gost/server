package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/odata"
	"github.com/stretchr/testify/assert"
)

func testHandlerGet() (*entities.Thing, error) {
	return nil, nil
}

func testHandlerGetError() (*entities.Thing, error) {
	return nil, errors.New("Test error")
}

func TestHandleGetTestWithQueroOptionsError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/things?$sort=bla&$skip='10'", nil)
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return testHandlerGet() }

	// act
	handleGetRequest(rr, nil, req, &handle, false, 10, "")

	// assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandleGetTestWithError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/things?$filter=name nonsens 'test1'", nil)
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return testHandlerGetError() }

	// act
	handleGetRequest(rr, nil, req, &handle, false, 10, "")

	// assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandleGetTestOk(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/bla", nil)
	req.Header.Set("Content-Type", "application/json")
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return testHandlerGet() }

	// act
	handleGetRequest(rr, nil, req, &handle, false, 10, "")

	// assert

	assert.Equal(t, http.StatusOK, rr.Code)
}
