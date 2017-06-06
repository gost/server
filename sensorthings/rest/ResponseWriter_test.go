package rest

import (
	"errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/odata"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendErrorWithNoError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()

	// act
	sendError(rr, nil)

	// assert
	assert.True(t, rr.Code == http.StatusInternalServerError)
}

func TestSendErrorWithNoIdentJson(t *testing.T) {
	// arrange
	IndentJSON = false
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	thing.Name = "yo"

	// act
	sendJSONResponse(rr, http.StatusTeapot, thing, nil)

	// assert
	assert.True(t, rr.Code == http.StatusTeapot)
}

func TestSendErrorWithError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	err1 := errors.New("wrong")
	errs := []error{err1}

	// act
	sendError(rr, errs)

	// assert
	assert.True(t, rr.Code == http.StatusInternalServerError)
}

func TestSendJsonResponseWithNoData(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()

	// act
	sendJSONResponse(rr, http.StatusTeapot, nil, nil)

	// assert
	assert.True(t, rr.Code == http.StatusTeapot)
}

func TestSendJsonResponseWithData(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	thing.Name = "yo"

	// act
	sendJSONResponse(rr, http.StatusTeapot, thing, nil)

	// assert
	assert.True(t, rr.Code == http.StatusTeapot)
}

func TestSendJsonResponseWithDataAndQueryOptions(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	thing.Name = "yo"
	req, _ := http.NewRequest("GET", "/v1.0/Things?$top=1&$select=name,id,description", nil)
	qo, _ := getQueryOptions(req)

	val := odata.GoDataValueQuery(true)
	qo.Value = &val

	// act
	sendJSONResponse(rr, http.StatusTeapot, thing, qo)

	// assert
	assert.True(t, rr.Code == http.StatusTeapot)
}
