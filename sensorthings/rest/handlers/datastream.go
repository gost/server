package handlers

import (
	"net/http"

	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/odata"
	"github.com/gost/server/sensorthings/rest/reader"
)

// HandleGetDatastreams retrieves datastreams based on Query Parameters
func HandleGetDatastreams(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetDatastreams(q, path) }
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetDatastream retrieves a datastream by given id
func HandleGetDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetDatastream(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetDatastreamByObservation ...
func HandleGetDatastreamByObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetDatastreamByObservation(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetDatastreamsByThing ...
func HandleGetDatastreamsByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetDatastreamsByThing(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetDatastreamsBySensor ...
func HandleGetDatastreamsBySensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetDatastreamsBySensor(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetDatastreamsByObservedProperty ...
func HandleGetDatastreamsByObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetDatastreamsByObservedProperty(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandlePostDatastream ...
func HandlePostDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ds := &entities.Datastream{}
	handle := func() (interface{}, []error) { return a.PostDatastream(ds) }
	handlePostRequest(w, endpoint, r, ds, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePostDatastreamByThing ...
func HandlePostDatastreamByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ds := &entities.Datastream{}
	handle := func() (interface{}, []error) { return a.PostDatastreamByThing(reader.GetEntityID(r), ds) }
	handlePostRequest(w, endpoint, r, ds, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandleDeleteDatastream ...
func HandleDeleteDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteDatastream(reader.GetEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePatchDatastream ...
func HandlePatchDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ds := &entities.Datastream{}
	handle := func() (interface{}, error) { return a.PatchDatastream(reader.GetEntityID(r), ds) }
	handlePatchRequest(w, endpoint, r, ds, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePutDatastream ...
func HandlePutDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ds := &entities.Datastream{}
	handle := func() (interface{}, []error) { return a.PutDatastream(reader.GetEntityID(r), ds) }
	handlePutRequest(w, endpoint, r, ds, &handle, a.GetConfig().Server.IndentedJSON)
}
