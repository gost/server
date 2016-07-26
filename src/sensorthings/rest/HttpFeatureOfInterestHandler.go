package rest

import (
	"net/http"

	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// HandleGetFeatureOfInterests ...
func HandleGetFeatureOfInterests(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetFeatureOfInterests(q, path) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetFeatureOfInterest ...
func HandleGetFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetFeatureOfInterest(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetFeatureOfInterestByObservation ...
func HandleGetFeatureOfInterestByObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetFeatureOfInterestByObservation(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostFeatureOfInterest ...
func HandlePostFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	foi := &entities.FeatureOfInterest{}
	handle := func() (interface{}, []error) { return a.PostFeatureOfInterest(foi) }
	handlePostRequest(w, endpoint, r, foi, &handle)
}

// HandleDeleteFeatureOfInterest ...
func HandleDeleteFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteFeatureOfInterest(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchFeatureOfInterest ...
func HandlePatchFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	foi := &entities.FeatureOfInterest{}
	handle := func() (interface{}, error) { return a.PatchFeatureOfInterest(getEntityID(r), foi) }
	handlePatchRequest(w, endpoint, r, foi, &handle)
}
