package handlers

import (
	"net/http"

	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/odata"
	"github.com/gost/server/sensorthings/rest/reader"
)

// HandleGetObservedProperty retrieves an ObservedProperty by id
func HandleGetObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservedProperty(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetObservedProperties retrieves ObservedProperties
func HandleGetObservedProperties(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetObservedProperties(q, path) }
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetObservedPropertyByDatastream retrieves the ObservedProperty by given Datastream id
func HandleGetObservedPropertyByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservedPropertyByDatastream(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandlePostObservedProperty posts a new ObservedProperty
func HandlePostObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	op := &entities.ObservedProperty{}
	handle := func() (interface{}, []error) { return a.PostObservedProperty(op) }
	handlePostRequest(w, endpoint, r, op, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandleDeleteObservedProperty Deletes an ObservedProperty by id
func HandleDeleteObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteObservedProperty(reader.GetEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePatchObservedProperty patches an Observes property by id
func HandlePatchObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	op := &entities.ObservedProperty{}
	handle := func() (interface{}, error) { return a.PatchObservedProperty(reader.GetEntityID(r), op) }
	handlePatchRequest(w, endpoint, r, op, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePutObservedProperty posts a new ObservedProperty
func HandlePutObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	op := &entities.ObservedProperty{}
	handle := func() (interface{}, []error) { return a.PutObservedProperty(reader.GetEntityID(r), op) }
	handlePutRequest(w, endpoint, r, op, &handle, a.GetConfig().Server.IndentedJSON)
}
