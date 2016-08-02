package rest

import (
	"net/http"

	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// HandleGetObservedProperty retrieves an ObservedProperty by id
func HandleGetObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservedProperty(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetObservedProperties retrieves ObservedProperties
func HandleGetObservedProperties(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetObservedProperties(q, path) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetObservedPropertyByDatastream retrieves the ObservedProperty by given Datastream id
func HandleGetObservedPropertyByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservedPropertyByDatastream(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostObservedProperty posts a new ObservedProperty
func HandlePostObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	op := &entities.ObservedProperty{}
	handle := func() (interface{}, []error) { return a.PostObservedProperty(op) }
	handlePostRequest(w, endpoint, r, op, &handle)
}

// HandleDeleteObservedProperty Deletes an ObservedProperty by id
func HandleDeleteObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteObservedProperty(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchObservedProperty patches an Observes property by id
func HandlePatchObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	op := &entities.ObservedProperty{}
	handle := func() (interface{}, error) { return a.PatchObservedProperty(getEntityID(r), op) }
	handlePatchRequest(w, endpoint, r, op, &handle)
}

// HandlePutObservedProperty posts a new ObservedProperty
func HandlePutObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	op := &entities.ObservedProperty{}
	handle := func() (interface{}, []error) { return a.PutObservedProperty(getEntityID(r), op) }
	handlePutRequest(w, endpoint, r, op, &handle)
}
