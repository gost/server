package http

import (
	"fmt"
	"net/http"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/gorilla/mux"
	"sort"
)

// CreateRouter creates a new mux.Router and sets up all endpoints defind in the sensothings api
func CreateRouter(api *models.API) *mux.Router {
	// Note: tried julienschmidt/httprouter instead of gorilla/mux but had some
	// problems with interfering endpoints cause of the wildcard used for the (id) in requests
	a := *api

	// get all endpoints into HttpEndpoints to be able to sort them so they can be added
	// to the routes in the right order else requests will be picked up by the wrong handlers
	eps := HttpEndpoints{}
	for _, endpoint := range *a.GetEndpoints() {
		for _, op := range endpoint.GetOperations() {
			e := &HttpEndpoint{Endpoint: endpoint, Operation: op}
			eps = append(eps, e)
		}
	}
	sort.Sort(eps)

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/Dashboard/").Handler(http.StripPrefix("/Dashboard/", http.FileServer(http.Dir(a.GetConfig().Server.ClientContent))))
	setDashboardRedirects(router)

	for _, e := range eps {
		op := e
		operation := op.Operation
		method := fmt.Sprintf("%s", operation.OperationType)
		if operation.Handler == nil {
			continue
		}

		router.Methods(method).
			Path(operation.Path).
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				operation.Handler(w, r, &op.Endpoint, api)
			})
	}

	return router
}

func setDashboardRedirects(router *mux.Router) {
	router.Methods("GET").Path("/Dashboard").HandlerFunc(dashboardRedirector)
	router.Methods("GET").Path("/dashboard").HandlerFunc(dashboardRedirector)
	router.Methods("GET").Path("/dashboard/").HandlerFunc(dashboardRedirector)
	router.Methods("GET").Path("").HandlerFunc(dashboardRedirector)
	router.Methods("GET").Path("/").HandlerFunc(dashboardRedirector)
}

func dashboardRedirector(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.URL.Host+"/Dashboard/", 302)
}
