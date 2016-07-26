package rest

import (
	"net/http"

	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// HandleGetSensorByDatastream ...
func HandleGetSensorByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetSensorByDatastream(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetSensor ...
func HandleGetSensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetSensor(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetSensors ...
func HandleGetSensors(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetSensors(q, path) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostSensors ...
func HandlePostSensors(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	sensor := &entities.Sensor{}
	handle := func() (interface{}, []error) { return a.PostSensor(sensor) }
	handlePostRequest(w, endpoint, r, sensor, &handle)
}

// HandleDeleteSensor ...
func HandleDeleteSensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteSensor(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchSensor ...
func HandlePatchSensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	sensor := &entities.Sensor{}
	handle := func() (interface{}, error) { return a.PatchSensor(getEntityID(r), sensor) }
	handlePatchRequest(w, endpoint, r, sensor, &handle)
}
