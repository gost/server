package http

import (
	"fmt"
	"net/http"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/gorilla/mux"
)

// CreateRouter creates a new mux.Router and sets up all endpoints defind in the sensothings api
func CreateRouter(api *models.API) *mux.Router {
	// Note: tried julienschmidt/httprouter instead of gorilla/mux but had some
	// problems with interfering endpoints cause of the wildcard used for the (id) in requests
	a := *api
	endpoints := *a.GetEndpoints()
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/Dashboard/").Handler(http.StripPrefix("/Dashboard/", http.FileServer(http.Dir("./clientv2/"))))
	setDashboardRedirects(router)

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
