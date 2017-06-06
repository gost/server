package rest

import (
	"net/http"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// HandleGetObservations ...
func HandleGetObservations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetObservations(q, path) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetObservation ...
func HandleGetObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservation(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetObservationsByFeatureOfInterest ...
func HandleGetObservationsByFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservationsByFeatureOfInterest(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetObservationsByDatastream ...
func HandleGetObservationsByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservationsByDatastream(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostObservation ...
func HandlePostObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.Observation{}
	handle := func() (interface{}, []error) { return a.PostObservation(ob) }
	handlePostRequest(w, endpoint, r, ob, &handle)
}

// HandlePostObservationByDatastream ...
func HandlePostObservationByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.Observation{}
	handle := func() (interface{}, []error) { return a.PostObservationByDatastream(getEntityID(r), ob) }
	handlePostRequest(w, endpoint, r, ob, &handle)
}

// HandleDeleteObservation ...
func HandleDeleteObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteObservation(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchObservation ...
func HandlePatchObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.Observation{}
	handle := func() (interface{}, error) { return a.PatchObservation(getEntityID(r), ob) }
	handlePatchRequest(w, endpoint, r, ob, &handle)
}

// HandlePutObservation ...
func HandlePutObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.Observation{}
	handle := func() (interface{}, []error) { return a.PutObservation(getEntityID(r), ob) }
	handlePutRequest(w, endpoint, r, ob, &handle)
}
