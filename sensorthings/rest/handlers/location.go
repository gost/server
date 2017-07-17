package handlers

import (
	"net/http"

	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/odata"
	"github.com/gost/server/sensorthings/rest/reader"
)

// HandleGetLocations retrieves multiple locations based on query parameters
func HandleGetLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetLocations(q, path) }
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetLocationsByHistoricalLocations retrieves the locations linked to the given Historical Location (id)
func HandleGetLocationsByHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetLocationsByHistoricalLocation(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetLocationsByThing retrieves the locations by given thing (id)
func HandleGetLocationsByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetLocationsByThing(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetLocation retrieves a location by given id
func HandleGetLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetLocation(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandlePostLocation posts a new location
func HandlePostLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, []error) { return a.PostLocation(loc) }
	handlePostRequest(w, endpoint, r, loc, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePostLocationByThing posts a new location linked to the given thing
func HandlePostLocationByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, []error) { return a.PostLocationByThing(reader.GetEntityID(r), loc) }
	handlePostRequest(w, endpoint, r, loc, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandleDeleteLocation deletes a location
func HandleDeleteLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteLocation(reader.GetEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePatchLocation patches a location by given id
func HandlePatchLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, error) { return a.PatchLocation(reader.GetEntityID(r), loc) }
	handlePatchRequest(w, endpoint, r, loc, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePutLocation patches a location by given id
func HandlePutLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, []error) { return a.PutLocation(reader.GetEntityID(r), loc) }
	handlePutRequest(w, endpoint, r, loc, &handle, a.GetConfig().Server.IndentedJSON)
}
