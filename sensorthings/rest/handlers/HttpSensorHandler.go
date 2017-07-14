package handlers

import (
	"net/http"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
	"github.com/geodan/gost/sensorthings/rest/reader"
)

// HandleGetSensorByDatastream ...
func HandleGetSensorByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetSensorByDatastream(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetSensor ...
func HandleGetSensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetSensor(reader.GetEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandleGetSensors ...
func HandleGetSensors(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetSensors(q, path) }
	handleGetRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON, a.GetConfig().Server.MaxEntityResponse, a.GetConfig().Server.ExternalURI)
}

// HandlePostSensors ...
func HandlePostSensors(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	sensor := &entities.Sensor{}
	handle := func() (interface{}, []error) { return a.PostSensor(sensor) }
	handlePostRequest(w, endpoint, r, sensor, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandleDeleteSensor ...
func HandleDeleteSensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteSensor(reader.GetEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePatchSensor ...
func HandlePatchSensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	sensor := &entities.Sensor{}
	handle := func() (interface{}, error) { return a.PatchSensor(reader.GetEntityID(r), sensor) }
	handlePatchRequest(w, endpoint, r, sensor, &handle, a.GetConfig().Server.IndentedJSON)
}

// HandlePutSensor ...
func HandlePutSensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	sensor := &entities.Sensor{}
	handle := func() (interface{}, []error) { return a.PutSensor(reader.GetEntityID(r), sensor) }
	handlePutRequest(w, endpoint, r, sensor, &handle, a.GetConfig().Server.IndentedJSON)
}
