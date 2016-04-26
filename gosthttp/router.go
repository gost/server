package gosthttp

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/geodan/gost/sensorthings"
)

func NewRouter(api *sensorthings.SensorThingsApi) *mux.Router {
	// Note: tried julienschmidt/httprouter instead of gorilla/mux but had some
	// problems with interfering endpoints cause of the wildcard used for the (id) in requests
	a := *api
	endpoints := *a.GetEndpoints()
	router := mux.NewRouter().StrictSlash(false)
	router.Path("/").Handler(http.FileServer(http.Dir("./gostsite/")))

	for _, endpoint := range endpoints {
		ep := endpoint
		for _, op := range ep.Operations {
			operation := op
			method := fmt.Sprintf("%s", operation.OperationType)
			if operation.Handler == nil {
				continue
			}

			router.Methods(method).
				Path(operation.Path).
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					operation.Handler(w, r, &ep, api)
				})
		}
	}

	return router
}
