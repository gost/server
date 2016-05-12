package rest

import (
	"encoding/json"
	"net/http"

	"io/ioutil"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
	"github.com/gorilla/mux"
)

// HandleAPIRoot will return a JSON array of the available SensorThings resource endpoints.
func HandleAPIRoot(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	bpi := a.GetBasePathInfo()
	sendJSONResponse(w, http.StatusOK, bpi)
}

// HandleVersion retrieves current version information and sends it back to the user
func HandleVersion(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	versionInfo := a.GetVersionInfo()
	sendJSONResponse(w, http.StatusOK, versionInfo)
}

// HandleGetThings retrieves and sends Things based on the given filter if provided
func HandleGetThings(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetThings(q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetThing retrieves and sends a specific Thing based on the given ID and filter
func HandleGetThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetThing(getEntityID(r), q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetThingByDatastream retrieves and sends a specific Thing based on the given datastream ID and filter
func HandleGetThingByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetThingByDatastream(getEntityID(r), q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostThing tries to insert a new Thing and sends back the created Thing
func HandlePostThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	thing := &entities.Thing{}
	handle := func() (interface{}, []error) { t := *thing; return a.PostThing(t) }
	handlePostRequest(w, endpoint, r, thing, &handle)
}

// HandleDeleteThing todo
func HandleDeleteThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteThing(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchThing todo
func HandlePatchThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	thing := &entities.Thing{}
	handle := func() (interface{}, error) { t := *thing; return a.PatchThing(getEntityID(r), t) }
	handlePatchRequest(w, endpoint, r, thing, &handle)
}

// HandleGetObservedProperty todo
func HandleGetObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetObservedProperty(getEntityID(r), q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetObservedProperties todo
func HandleGetObservedProperties(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetObservedProperties(q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetObservedPropertyByDatastream todo
func HandleGetObservedPropertyByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) {
		return a.GetObservedPropertiesByDatastream(getEntityID(r), q)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostObservedProperty todo
func HandlePostObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	op := &entities.ObservedProperty{}
	handle := func() (interface{}, []error) { o := *op; return a.PostObservedProperty(o) }
	handlePostRequest(w, endpoint, r, op, &handle)
}

// HandleDeleteObservedProperty todo
func HandleDeleteObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteObservedProperty(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchObservedProperty todo
func HandlePatchObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	op := &entities.ObservedProperty{}
	handle := func() (interface{}, error) { o := *op; return a.PatchObservedProperty(getEntityID(r), o) }
	handlePatchRequest(w, endpoint, r, op, &handle)
}

// HandleGetLocations todo
func HandleGetLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetLocations(q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetLocationsByThing todo
func HandleGetLocationsByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetLocationsByThing(getEntityID(r), q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetLocation todo
func HandleGetLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetLocation(getEntityID(r), q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostLocation todo
func HandlePostLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, []error) { l := *loc; return a.PostLocation(l) }
	handlePostRequest(w, endpoint, r, loc, &handle)
}

// HandlePostLocationByThing todo
func HandlePostLocationByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, []error) { l := *loc; return a.PostLocationByThing(getEntityID(r), l) }
	handlePostRequest(w, endpoint, r, loc, &handle)
}

// HandleDeleteLocation todo
func HandleDeleteLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteLocation(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchLocation todo
func HandlePatchLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, error) { l := *loc; return a.PatchLocation(getEntityID(r), l) }
	handlePatchRequest(w, endpoint, r, loc, &handle)
}

// HandleGetDatastreams todo
func HandleGetDatastreams(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetDatastreams(q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetDatastream todo
func HandleGetDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetDatastream(getEntityID(r), q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetDatastreamsByThing todo
func HandleGetDatastreamsByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetDatastreamsByThing(getEntityID(r), q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetDatastreamsBySensor todo
func HandleGetDatastreamsBySensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetDatastreamsBySensor(getEntityID(r), q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostDatastream todo
func HandlePostDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ds := &entities.Datastream{}
	handle := func() (interface{}, []error) { d := *ds; return a.PostDatastream(d) }
	handlePostRequest(w, endpoint, r, ds, &handle)
}

// HandlePostDatastreamByThing todo
func HandlePostDatastreamByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ds := &entities.Datastream{}
	handle := func() (interface{}, []error) { d := *ds; return a.PostDatastreamByThing(getEntityID(r), d) }
	handlePostRequest(w, endpoint, r, ds, &handle)
}

// HandleDeleteDatastream todo
func HandleDeleteDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteDatastream(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchDatastream todo
func HandlePatchDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ds := &entities.Datastream{}
	handle := func() (interface{}, error) { d := *ds; return a.PatchDatastream(getEntityID(r), d) }
	handlePatchRequest(w, endpoint, r, ds, &handle)
}

// HandleGetSensor todo
func HandleGetSensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetSensor(getEntityID(r), q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetSensors todo
func HandleGetSensors(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetSensors(q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostSensors todo
func HandlePostSensors(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	sensor := &entities.Sensor{}
	handle := func() (interface{}, []error) { s := *sensor; return a.PostSensor(s) }
	handlePostRequest(w, endpoint, r, sensor, &handle)
}

// HandleDeleteSensor todo
func HandleDeleteSensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteSensor(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchSensor todo
func HandlePatchSensor(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	sensor := &entities.Sensor{}
	handle := func() (interface{}, error) { s := *sensor; return a.PatchSensor(getEntityID(r), s) }
	handlePatchRequest(w, endpoint, r, sensor, &handle)
}

// HandleGetObservations todo
func HandleGetObservations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetObservations(q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetObservation todo
func HandleGetObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetObservation(getEntityID(r), q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetObservationsByDatastream todo
func HandleGetObservationsByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) {
		return a.GetObservationsByDatastream(getEntityID(r), q)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostObservation todo
func HandlePostObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.Observation{}
	handle := func() (interface{}, []error) { o := *ob; return a.PostObservation(o) }
	handlePostRequest(w, endpoint, r, ob, &handle)
}

// HandlePostObservationByDatastream todo
func HandlePostObservationByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.Observation{}
	handle := func() (interface{}, []error) { o := *ob; return a.PostObservationByDatastream(getEntityID(r), o) }
	handlePostRequest(w, endpoint, r, ob, &handle)
}

// HandleDeleteObservation todo
func HandleDeleteObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteObservation(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchObservation todo
func HandlePatchObservation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.Observation{}
	handle := func() (interface{}, error) { o := *ob; return a.PatchObservation(getEntityID(r), o) }
	handlePatchRequest(w, endpoint, r, ob, &handle)
}

// HandleGetFeatureOfInterests todo
func HandleGetFeatureOfInterests(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetFeatureOfInterests(q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetFeatureOfInterest todo
func HandleGetFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetFeatureOfInterest(getEntityID(r), q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostFeatureOfInterest todo
func HandlePostFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	foi := &entities.FeatureOfInterest{}
	handle := func() (interface{}, []error) { f := *foi; return a.PostFeatureOfInterest(f) }
	handlePostRequest(w, endpoint, r, foi, &handle)
}

// HandleDeleteFeatureOfInterest todo
func HandleDeleteFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteFeatureOfInterest(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchFeatureOfInterest todo
func HandlePatchFeatureOfInterest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	foi := &entities.FeatureOfInterest{}
	handle := func() (interface{}, error) { f := *foi; return a.PatchFeatureOfInterest(getEntityID(r), f) }
	handlePatchRequest(w, endpoint, r, foi, &handle)
}

// HandleGetHistoricalLocations todo
func HandleGetHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetHistoricalLocations(q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetHistoricalLocationsByThing todo
func HandleGetHistoricalLocationsByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions) (interface{}, error) {
		return a.GetHistoricalLocationsByThing(getEntityID(r), q)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetHistoricalLocation todo
func HandleGetHistoricalLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	id := getEntityID(r)
	handle := func(q *odata.QueryOptions) (interface{}, error) { return a.GetHistoricalLocation(id, q) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleDeleteHistoricalLocations todo
func HandleDeleteHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteHistoricalLocation(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchHistoricalLocations todo
func HandlePatchHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	hl := &entities.HistoricalLocation{}
	handle := func() (interface{}, error) { h := *hl; return a.PatchHistoricalLocation(getEntityID(r), h) }
	handlePatchRequest(w, endpoint, r, hl, &handle)
}

// getEntityID retrieves the id from the request, for example
// the request http://mysensor.com/V1.0/Things(1236532) returns 1236532 as id
func getEntityID(r *http.Request) string {
	vars := mux.Vars(r)
	value := vars["id"]
	substring := value[1 : len(value)-1]
	return substring
}

// GetQueryOptions creates QueryOptions based upon the incoming request
// QueryOptions = nil when no query was found, errors != nil if something
// went wrong with parsing the query into QueryOptions and will contain information
// on what went wrong
func getQueryOptions(r *http.Request) (*odata.QueryOptions, []error) {
	query := r.URL.Query()
	if len(query) == 0 {
		return nil, nil
	}

	qo, e := odata.CreateQueryOptions(query)
	return qo, e
}

// handleGetRequest is the default function to handle incoming GET requests
func handleGetRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, h *func(q *odata.QueryOptions) (interface{}, error)) {
	// Parse query options from request
	queryOptions, err := getQueryOptions(r)
	if err != nil {
		sendError(w, err)
		return
	}

	// Check if the requested enpoints supports the parsed queries
	endpoint := *e
	_, err = endpoint.AreQueryOptionsSupported(queryOptions)
	if err != nil {
		sendError(w, err)
		return
	}

	// Run the handler func such as Api.GetThingById
	handler := *h
	data, err2 := handler(queryOptions)
	if err2 != nil {
		sendError(w, []error{err2})
		return
	}

	sendJSONResponse(w, http.StatusOK, data)
}

// handlePatchRequest todo: currently almost same as handlePostRequest, merge if it stays like this
func handlePatchRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, entity entities.Entity, h *func() (interface{}, error)) {
	byteData, _ := ioutil.ReadAll(r.Body)
	err := entity.ParseEntity(byteData)
	if err != nil {
		sendError(w, []error{err})
		return
	}

	handle := *h
	data, err2 := handle()
	if err2 != nil {
		sendError(w, []error{err2})
		return
	}

	sendJSONResponse(w, http.StatusOK, data)
}

// handlePostRequest todo
func handleDeleteRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, h *func() error) {
	handle := *h
	err := handle()
	if err != nil {
		sendError(w, []error{err})
		return
	}

	sendJSONResponse(w, http.StatusOK, "null")
}

// handlePostRequest todo
func handlePostRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, entity entities.Entity, h *func() (interface{}, []error)) {
	byteData, _ := ioutil.ReadAll(r.Body)
	err := entity.ParseEntity(byteData)
	if err != nil {
		sendError(w, []error{err})
		return
	}

	handle := *h
	data, err2 := handle()
	if err2 != nil {
		sendError(w, err2)
		return
	}

	sendJSONResponse(w, http.StatusCreated, data)
}

// sendJSONResponse sends the desired message to the user
// the message will be marshalled into an indented JSON format
func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	b, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		panic(err)
	}

	w.Write(b)
}

// sendError creates an ErrorResponse message and sets it to the user
// using SendJSONResponse
func sendError(w http.ResponseWriter, error []error) {
	//errors cannot be marshalled, create strings
	errors := make([]string, len(error))
	for idx, value := range error {
		errors[idx] = value.Error()
	}

	// Set te status code, default 500 for error, check if there is an ApiError an get
	// the status code
	var statusCode = http.StatusInternalServerError
	if error != nil && len(error) > 0 {
		switch e := error[0].(type) {
		case gostErrors.APIError:
			statusCode = e.GetHTTPErrorStatusCode()
			break
		}
	}

	statusText := http.StatusText(statusCode)
	errorResponse := models.ErrorResponse{
		Error: models.ErrorContent{
			StatusText: statusText,
			StatusCode: statusCode,
			Messages:   errors,
		},
	}

	sendJSONResponse(w, statusCode, errorResponse)
}
