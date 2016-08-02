package rest

import (
	"net/http"

	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// HandleGetHistoricalLocations ...
func HandleGetHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocations(q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetHistoricalLocationsByThing ...
func HandleGetHistoricalLocationsByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocationsByThing(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetHistoricalLocationsByLocation ...
func HandleGetHistoricalLocationsByLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocationsByLocation(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetHistoricalLocation ...
func HandleGetHistoricalLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	id := getEntityID(r)
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocation(id, q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostHistoricalLocation ...
func HandlePostHistoricalLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	hl := &entities.HistoricalLocation{}
	handle := func() (interface{}, []error) { return a.PostHistoricalLocation(hl) }
	handlePostRequest(w, endpoint, r, hl, &handle)
}

// HandlePutHistoricalLocation ...
func HandlePutHistoricalLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	hl := &entities.HistoricalLocation{}
	handle := func() (interface{}, []error) { return a.PutHistoricalLocation(getEntityID(r), hl) }
	handlePutRequest(w, endpoint, r, hl, &handle)
}

// HandleDeleteHistoricalLocations ...
func HandleDeleteHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteHistoricalLocation(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchHistoricalLocations ...
func HandlePatchHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	hl := &entities.HistoricalLocation{}
	handle := func() (interface{}, error) { return a.PatchHistoricalLocation(getEntityID(r), hl) }
	handlePatchRequest(w, endpoint, r, hl, &handle)
}
