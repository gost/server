package rest

import (
	"net/http"

	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// HandleGetDatastreams retrieves datastreams based on Query Parameters
func HandleGetDatastreams(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetDatastreams(q, path) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetDatastream retrieves a datastream by given id
func HandleGetDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetDatastream(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetDatastreamByObservation ...
func HandleGetDatastreamByObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetDatastreamByObservation(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetDatastreamsByThing ...
func HandleGetDatastreamsByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetDatastreamsByThing(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetDatastreamsBySensor ...
func HandleGetDatastreamsBySensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetDatastreamsBySensor(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetDatastreamsByObservedProperty ...
func HandleGetDatastreamsByObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetDatastreamsByObservedProperty(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostDatastream ...
func HandlePostDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ds := &entities.Datastream{}
	handle := func() (interface{}, []error) { return a.PostDatastream(ds) }
	handlePostRequest(w, endpoint, r, ds, &handle)
}

// HandlePostDatastreamByThing ...
func HandlePostDatastreamByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ds := &entities.Datastream{}
	handle := func() (interface{}, []error) { return a.PostDatastreamByThing(getEntityID(r), ds) }
	handlePostRequest(w, endpoint, r, ds, &handle)
}

// HandleDeleteDatastream ...
func HandleDeleteDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteDatastream(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchDatastream ...
func HandlePatchDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ds := &entities.Datastream{}
	handle := func() (interface{}, error) { return a.PatchDatastream(getEntityID(r), ds) }
	handlePatchRequest(w, endpoint, r, ds, &handle)
}
