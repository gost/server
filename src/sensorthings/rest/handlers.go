package rest

import (
	"encoding/json"
	"net/http"

	"io/ioutil"

	"bytes"
	"fmt"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
	"github.com/gorilla/mux"
	"strings"
)

// HandleAPIRoot will return a JSON array of the available SensorThings resource endpoints.
func HandleAPIRoot(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	bpi := a.GetBasePathInfo()
	sendJSONResponse(w, http.StatusOK, bpi, nil)
}

// HandleVersion retrieves current version information and sends it back to the user
func HandleVersion(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	versionInfo := a.GetVersionInfo()
	sendJSONResponse(w, http.StatusOK, versionInfo, nil)
}

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

// HandleGetObservedProperty retrieves an ObservedProperty by id
func HandleGetObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservedProperty(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetObservedProperties retrieves ObservedProperties
func HandleGetObservedProperties(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetObservedProperties(q, path) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetObservedPropertyByDatastream retrieves the ObservedProperty by given Datastream id
func HandleGetObservedPropertyByDatastream(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetObservedPropertyByDatastream(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostObservedProperty posts a new ObservedProperty
func HandlePostObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	op := &entities.ObservedProperty{}
	handle := func() (interface{}, []error) { return a.PostObservedProperty(op) }
	handlePostRequest(w, endpoint, r, op, &handle)
}

// HandleDeleteObservedProperty Deletes an ObservedProperty by id
func HandleDeleteObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteObservedProperty(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchObservedProperty patches an Observes property by id
func HandlePatchObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	op := &entities.ObservedProperty{}
	handle := func() (interface{}, error) { return a.PatchObservedProperty(getEntityID(r), op) }
	handlePatchRequest(w, endpoint, r, op, &handle)
}

// HandleGetLocations retrieves multiple locations based on query parameters
func HandleGetLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) { return a.GetLocations(q, path) }
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetLocationsByHistoricalLocations retrieves the locations linked to the given Historical Location (id)
func HandleGetLocationsByHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetLocationsByHistoricalLocation(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetLocationsByThing retrieves the locations by given thing (id)
func HandleGetLocationsByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetLocationsByThing(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetLocation retrieves a location by given id
func HandleGetLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetLocation(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandlePostLocation posts a new location
func HandlePostLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, []error) { return a.PostLocation(loc) }
	handlePostRequest(w, endpoint, r, loc, &handle)
}

// HandlePostLocationByThing posts a new location linked to the given thing
func HandlePostLocationByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, []error) { return a.PostLocationByThing(getEntityID(r), loc) }
	handlePostRequest(w, endpoint, r, loc, &handle)
}

// HandleDeleteLocation deletes a location
func HandleDeleteLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteLocation(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchLocation patches a location by given id
func HandlePatchLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	loc := &entities.Location{}
	handle := func() (interface{}, error) { return a.PatchLocation(getEntityID(r), loc) }
	handlePatchRequest(w, endpoint, r, loc, &handle)
}

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
	handle := func() (interface{}, error) { return a.PatchSensor(getEntityID(r), &entities.Sensor{}) }
	handlePatchRequest(w, endpoint, r, sensor, &handle)
}

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

// HandleGetHistoricalLocations ...
func HandleGetHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocations(q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetHistoricalLocationsByThing ...
func HandleGetHistoricalLocationsByThing(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocationsByThing(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetHistoricalLocationsByLocation ...
func HandleGetHistoricalLocationsByLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocationsByLocation(getEntityID(r), q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleGetHistoricalLocation ...
func HandleGetHistoricalLocation(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	id := getEntityID(r)
	handle := func(q *odata.QueryOptions, path string) (interface{}, error) {
		return a.GetHistoricalLocation(id, q, path)
	}
	handleGetRequest(w, endpoint, r, &handle)
}

// HandleDeleteHistoricalLocations ...
func HandleDeleteHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	handle := func() error { return a.DeleteHistoricalLocation(getEntityID(r)) }
	handleDeleteRequest(w, endpoint, r, &handle)
}

// HandlePatchHistoricalLocations ...
func HandlePatchHistoricalLocations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	hl := &entities.HistoricalLocation{}
	handle := func() (interface{}, error) { return a.PatchHistoricalLocation(getEntityID(r), hl) }
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

	//If request contains parameters from route wildcard convert it to a select query
	vars := mux.Vars(r)
	value := []string{vars["params"]}

	if len(vars["params"]) > 0 {
		//If $ref found create select query with id
		if vars["params"] == "$ref" {
			value = []string{"id"}
			query["$ref"] = []string{"true"}
		}

		query["$select"] = value
	}

	if strings.HasSuffix(r.URL.Path, "$value") {
		query["$value"] = []string{"true"}
	}

	// if $top is not found, retrieve max 200
	_, ok := query["$top"]
	if !ok {
		query["$top"] = []string{"200"}
	}

	_, ok = query["$skip"]
	if !ok {
		query["$skip"] = []string{"0"}
	}

	qo, e := odata.CreateQueryOptions(query)
	return qo, e
}

// handleGetRequest is the default function to handle incoming GET requests
func handleGetRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, h *func(q *odata.QueryOptions, path string) (interface{}, error)) {
	// Parse query options from request
	queryOptions, err := getQueryOptions(r)
	if err != nil {
		sendError(w, err)
		return
	}

	// Check if the requested endpoints supports the parsed queries
	endpoint := *e
	_, err = endpoint.AreQueryOptionsSupported(queryOptions)
	if err != nil {
		sendError(w, err)
		return
	}

	// Run the handler func such as Api.GetThingById
	handler := *h
	data, err2 := handler(queryOptions, fmt.Sprintf(r.Host+r.URL.Path))
	if err2 != nil {
		sendError(w, []error{err2})
		return
	}

	sendJSONResponse(w, http.StatusOK, data, queryOptions)
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

	sendJSONResponse(w, http.StatusOK, data, nil)
}

// handlePostRequest
func handleDeleteRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, h *func() error) {
	handle := *h
	err := handle()
	if err != nil {
		sendError(w, []error{err})
		return
	}

	sendJSONResponse(w, http.StatusOK, "null", nil)
}

// handlePostRequest
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

	w.Header().Add("Location", entity.GetSelfLink())

	sendJSONResponse(w, http.StatusCreated, data, nil)
}

// sendJSONResponse sends the desired message to the user
// the message will be marshalled into an indented JSON format
func sendJSONResponse(w http.ResponseWriter, status int, data interface{}, qo *odata.QueryOptions) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	b, err := JSONMarshal(data, true)
	if err != nil {
		panic(err)
	}

	// $value is requested only send back the value, ToDo: move to API code?
	if qo != nil && qo.QueryOptionValue {
		errMessage := fmt.Errorf("Unable to retrieve $value for %v", qo.QuerySelect.Params[0])
		var m map[string]json.RawMessage
		err = json.Unmarshal(b, &m)
		if err != nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
			sendError(w, []error{gostErrors.NewRequestInternalServerError(errMessage)})
		}

		mVal, ok := m[qo.QuerySelect.Params[0]]
		if !ok {
			sendError(w, []error{gostErrors.NewRequestInternalServerError(errMessage)})
		}

		value := string(mVal[:])
		value = strings.TrimPrefix(value, "\"")
		value = strings.TrimSuffix(value, "\"")

		b = []byte(value)
	}

	w.Write(b)
}

//JSONMarshal converts the data and converts special characters such as &
func JSONMarshal(data interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.MarshalIndent(data, "", "   ")

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
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

	sendJSONResponse(w, statusCode, errorResponse, nil)
}
