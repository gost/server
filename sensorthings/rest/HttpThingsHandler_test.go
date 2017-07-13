package rest

import (
	"encoding/json"
	"fmt"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getRouter() *mux.Router {
	a := NewMockAPI()
	eps := EndpointsToSortedList(a.GetEndpoints())
	router := mux.NewRouter().StrictSlash(false)

	for _, e := range eps {
		op := e
		operation := op.Operation
		method := fmt.Sprintf("%s", operation.OperationType)
		router.Methods(method).
			Path(operation.Path).
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				operation.Handler(w, r, &op.Endpoint, &a)
			})
	}

	return router
}

func TestGetThing(t *testing.T) {
	// arrange
	mockThing := NewMockThing()
	router := getRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	// act
	r, _ := http.Get(ts.URL + "/v1.0/things(1)")
	thing := entities.Thing{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &thing)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, fmt.Sprintf("%v", mockThing.ID), fmt.Sprintf("%v", thing.ID))
	assert.Equal(t, mockThing.Name, thing.Name)
	assert.Equal(t, mockThing.Description, thing.Description)
	assert.Equal(t, mockThing.Properties, thing.Properties)
}
