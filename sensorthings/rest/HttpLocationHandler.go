package rest

import (
	"net/http"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// HandleGetLocations retrieves multiple locations based on query parameters
func HandleGetLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetLocations(q, path) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetLocationsByHistoricalLocations retrieves the locations linked to the given Historical Location (id)
func HandleGetLocationsByHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetLocationsByHistoricalLocation(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetLocationsByThing retrieves the locations by given thing (id)
func HandleGetLocationsByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetLocationsByThing(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetLocation retrieves a location by given id
func HandleGetLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetLocation(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostLocation posts a new location
func HandlePostLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, []error) { return a.PostLocation(loc) }
	handlePostRequest(w, endpoint, r, loc, &handle)
}

// HandlePostLocationByThing posts a new location linked to the given thing
func HandlePostLocationByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, []error) { return a.PostLocationByThing(getEntityID(r), loc) }
	handlePostRequest(w, endpoint, r, loc, &handle)
}

// HandleDeleteLocation deletes a location
func HandleDeleteLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteLocation(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchLocation patches a location by given id
func HandlePatchLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, error) { return a.PatchLocation(getEntityID(r), loc) }
	handlePatchRequest(w, endpoint, r, loc, &handle)
}

// HandlePutLocation patches a location by given id
func HandlePutLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, []error) { return a.PutLocation(getEntityID(r), loc) }
	handlePutRequest(w, endpoint, r, loc, &handle)
}
