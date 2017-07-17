package handlers

import (
	"net/http"

	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/odata"
	"github.com/gost/server/sensorthings/rest/reader"
)

// HandleGetObservations ...
func HandleGetObservations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetObservations(q, path) }
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetObservation ...
func HandleGetObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservation(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetObservationsByFeatureOfInterest ...
func HandleGetObservationsByFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservationsByFeatureOfInterest(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetObservationsByDatastream ...
func HandleGetObservationsByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservationsByDatastream(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandlePostObservation ...
func HandlePostObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.Observation{}
	handle := func() (interface{}, []error) { return a.PostObservation(ob) }
	handlePostRequest(w, endpoint, r, ob, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePostObservationByDatastream ...
func HandlePostObservationByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.Observation{}
	handle := func() (interface{}, []error) { return a.PostObservationByDatastream(reader.GetEntityID(r), ob) }
	handlePostRequest(w, endpoint, r, ob, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandleDeleteObservation ...
func HandleDeleteObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteObservation(reader.GetEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePatchObservation ...
func HandlePatchObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.Observation{}
	handle := func() (interface{}, error) { return a.PatchObservation(reader.GetEntityID(r), ob) }
	handlePatchRequest(w, endpoint, r, ob, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePutObservation ...
func HandlePutObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.Observation{}
	handle := func() (interface{}, []error) { return a.PutObservation(reader.GetEntityID(r), ob) }
	handlePutRequest(w, endpoint, r, ob, &handle, a.GetConfig().Server.IndentedJSON)
}
