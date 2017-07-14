package handlers

import (
	"net/http"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
	"github.com/geodan/gost/sensorthings/rest/reader"
)

// HandleGetFeatureOfInterests ...
func HandleGetFeatureOfInterests(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetFeatureOfInterests(q, path) }
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetFeatureOfInterest ...
func HandleGetFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetFeatureOfInterest(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetFeatureOfInterestByObservation ...
func HandleGetFeatureOfInterestByObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetFeatureOfInterestByObservation(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandlePostFeatureOfInterest ...
func HandlePostFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	foi := &entities.FeatureOfInterest{}
	handle := func() (interface{}, []error) { return a.PostFeatureOfInterest(foi) }
	handlePostRequest(w, endpoint, r, foi, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandleDeleteFeatureOfInterest ...
func HandleDeleteFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteFeatureOfInterest(reader.GetEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePatchFeatureOfInterest ...
func HandlePatchFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	foi := &entities.FeatureOfInterest{}
	handle := func() (interface{}, error) { return a.PatchFeatureOfInterest(reader.GetEntityID(r), foi) }
	handlePatchRequest(w, endpoint, r, foi, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePutFeatureOfInterest ...
func HandlePutFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	foi := &entities.FeatureOfInterest{}
	handle := func() (interface{}, []error) { return a.PutFeatureOfInterest(reader.GetEntityID(r), foi) }
	handlePutRequest(w, endpoint, r, foi, &handle, a.GetConfig().Server.IndentedJSON)
}
