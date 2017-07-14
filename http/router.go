package http

import (
	"fmt"
	"net/http"

	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/rest/endpoint"
	"github.com/gorilla/mux"
)

// CreateRouter creates a new mux.Router and sets up all endpoints defined in the SensorThings api
func CreateRouter(api *models.API) *mux.Router {
	// Note: tried julienschmidt/httprouter instead of gorilla/mux but had some
	// problems with interfering endpoints cause of the wildcard used for the (id) in requests
	a := *api
	eps := endpoint.EndpointsToSortedList(a.GetEndpoints())
	router := mux.NewRouter().StrictSlash(false)

	for _, e := range eps {
		op := e
		operation := op.Operation
		method := fmt.Sprintf("%s", operation.OperationType)
		router.Methods(method).
			Path(operation.Path).
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				operation.Handler(w, r, &op.Endpoint, api)
			})
	}

	return router
}
