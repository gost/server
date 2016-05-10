package http

import (
	"fmt"
	"net/http"

	"github.com/geodan/gost/sensorthings/models"
	"github.com/gorilla/mux"
)

// CreateRouter creates a new mux.Router and sets up all endpoints defind in the sensothings api
func CreateRouter(api *models.API) *mux.Router {
	// Note: tried julienschmidt/httprouter instead of gorilla/mux but had some
	// problems with interfering endpoints cause of the wildcard used for the (id) in requests
	a := *api
	endpoints := *a.GetEndpoints()
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/dashboard").Handler(http.StripPrefix("/dashboard", http.FileServer(http.Dir("./client/"))))

	for _, endpoint := range endpoints {
		ep := endpoint
		for _, op := range ep.GetOperations() {
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
