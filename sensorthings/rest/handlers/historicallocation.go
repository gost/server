package handlers

import (
	"net/http"

	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/odata"
	"github.com/gost/server/sensorthings/rest/reader"
)

// HandleGetHistoricalLocations ...
func HandleGetHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocations(q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetHistoricalLocationsByThing ...
func HandleGetHistoricalLocationsByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocationsByThing(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetHistoricalLocationsByLocation ...
func HandleGetHistoricalLocationsByLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocationsByLocation(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetHistoricalLocation ...
func HandleGetHistoricalLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	id := reader.GetEntityID(r)
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocation(id, q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandlePostHistoricalLocation ...
func HandlePostHistoricalLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	hl := &entities.HistoricalLocation{}
	handle := func() (interface{}, []error) { return a.PostHistoricalLocation(hl) }
	handlePostRequest(w, endpoint, r, hl, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePutHistoricalLocation ...
func HandlePutHistoricalLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	hl := &entities.HistoricalLocation{}
	handle := func() (interface{}, []error) { return a.PutHistoricalLocation(reader.GetEntityID(r), hl) }
	handlePutRequest(w, endpoint, r, hl, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandleDeleteHistoricalLocations ...
func HandleDeleteHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteHistoricalLocation(reader.GetEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePatchHistoricalLocations ...
func HandlePatchHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	hl := &entities.HistoricalLocation{}
	handle := func() (interface{}, error) { return a.PatchHistoricalLocation(reader.GetEntityID(r), hl) }
	handlePatchRequest(w, endpoint, r, hl, &handle, a.GetConfig().Server.IndentedJSON)
}
