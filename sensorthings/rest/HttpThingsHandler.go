package rest

import (
	"net/http"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// HandleGetThings retrieves and sends Things based on the given filter if provided
func HandleGetThings(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetThings(q, path) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetThing retrieves and sends a specific Thing based on the given ID and filter
func HandleGetThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetThing(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetThingByDatastream retrieves and sends a specific Thing based on the given datastream ID and filter
func HandleGetThingByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetThingByDatastream(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetThingsByLocation retrieves and sends Things based on the given Location ID and filter
func HandleGetThingsByLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetThingsByLocation(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetThingByHistoricalLocation retrieves and sends a specific Thing based on the given HistoricalLocation ID and filter
func HandleGetThingByHistoricalLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetThingByHistoricalLocation(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostThing tries to insert a new Thing and sends back the created Thing
func HandlePostThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	thing := &entities.Thing{}
	handle := func() (interface{}, []error) { return a.PostThing(thing) }
	handlePostRequest(w, endpoint, r, thing, &handle)
}

// HandleDeleteThing deletes a thing by given id
func HandleDeleteThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteThing(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchThing patches a thing by given id
func HandlePatchThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	thing := &entities.Thing{}
	handle := func() (interface{}, error) { return a.PatchThing(getEntityID(r), thing) }
	handlePatchRequest(w, endpoint, r, thing, &handle)
}

// HandlePutThing patches a thing by given id
func HandlePutThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	thing := &entities.Thing{}
	handle := func() (interface{}, []error) { return a.PutThing(getEntityID(r), thing) }
	handlePutRequest(w, endpoint, r, thing, &handle)
}
