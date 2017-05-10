package rest

import (
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"errors"
	"github.com/geodan/gost/src/sensorthings/entities"

)

func TestSendErrorWithNoError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()

	// act
	sendError(rr, nil)

	// assert
	assert.True(t, rr.Code == http.StatusInternalServerError)
}

func TestSendErrorWithError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	err1 := errors.New("wrong")
	errs:= []error{err1}

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