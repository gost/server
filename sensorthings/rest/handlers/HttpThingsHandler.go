package handlers

import (
	"net/http"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
	"github.com/geodan/gost/sensorthings/rest/reader"
)

// HandleGetThings retrieves and sends Things based on the given filter if provided
func HandleGetThings(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetThings(q, path) }
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetThing retrieves and sends a specific Thing based on the given ID and filter
func HandleGetThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetThing(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetThingByDatastream retrieves and sends a specific Thing based on the given datastream ID and filter
func HandleGetThingByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetThingByDatastream(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetThingsByLocation retrieves and sends Things based on the given Location ID and filter
func HandleGetThingsByLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetThingsByLocation(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetThingByHistoricalLocation retrieves and sends a specific Thing based on the given HistoricalLocation ID and filter
func HandleGetThingByHistoricalLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetThingByHistoricalLocation(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandlePostThing tries to insert a new Thing and sends back the created Thing
func HandlePostThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	thing := &entities.Thing{}
	handle := func() (interface{}, []error) { return a.PostThing(thing) }
	handlePostRequest(w, endpoint, r, thing, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandleDeleteThing deletes a thing by given id
func HandleDeleteThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteThing(reader.GetEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePatchThing patches a thing by given id
func HandlePatchThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	thing := &entities.Thing{}
	handle := func() (interface{}, error) { return a.PatchThing(reader.GetEntityID(r), thing) }
	handlePatchRequest(w, endpoint, r, thing, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePutThing patches a thing by given id
func HandlePutThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	thing := &entities.Thing{}
	handle := func() (interface{}, []error) { return a.PutThing(reader.GetEntityID(r), thing) }
	handlePutRequest(w, endpoint, r, thing, &handle, a.GetConfig().Server.IndentedJSON)
}
