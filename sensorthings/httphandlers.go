package sensorthings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"io"
	"time"

	"github.com/gorilla/mux"
)

// HandleAPIRoot will return a JSON array of the available SensorThings resource endpoints.
func HandleAPIRoot(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
	a := *api
	bpi := a.GetBasePathInfo()
	sendJSONResponse(w, http.StatusOK, bpi)
}

// HandleVersion retrieves current version information and sends it back to the user
func HandleVersion(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
	a := *api
	versionInfo := a.GetVersionInfo()
	sendJSONResponse(w, http.StatusOK, versionInfo)
}

// HandleGetThings retrieves and sends Things based on the given filter if provided
func HandleGetThings(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
	fmt.Println(time.Now())
	a := *api
	handle := func(q *QueryOptions) (interface{}, error) { return a.GetThings(q) }
	handleGetRequest(w, endpoint, r, &handle, http.StatusOK)
	fmt.Println(time.Now())
}

// HandleGetThingByID retrieves and sends a specific Thing based on the given ID and filter
func HandleGetThingByID(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
	a := *api
	id := getEntityID(r)
	handle := func(q *QueryOptions) (interface{}, error) { return a.GetThing(id, q) }
	handleGetRequest(w, endpoint, r, &handle, http.StatusOK)
}

// HandlePostThing tries to insert a new Thing and sends back the created Thing with http.StatusCreated when successful
func HandlePostThing(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
	var t Thing
	if !tryParseEntity(w, r.Body, &t) {
		return
	}

	a := *api
	nt, err := a.PostThing(t)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
	} else {
		sendJSONResponse(w, http.StatusCreated, nt)
	}

	// test
	/*a := *api
	//var t Thing
	handle := func(test interface{}) (interface{}, []error) {
		asdf := *test
		return a.PostThing(asdf.(Thing))
	}

	handlePostRequest(w, endpoint, r, &handle, Thing{}, http.StatusOK)*/
}

// handlePostRequest todo
func handlePostRequest(w http.ResponseWriter, e *Endpoint, r *http.Request, h *func(interface{}) (interface{}, []error), i interface{}, statusCode int) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&i)
	if err != nil {
		sendError(w, http.StatusBadRequest, []error{err})
		return
	}

	handle := *h
	data, err2 := handle(&i)
	if err2 != nil {
		sendError(w, http.StatusBadRequest, []error{err})
		return
	}

	sendJSONResponse(w, statusCode, data)
}

// HandleDeleteThing todo
func HandleDeleteThing(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {

}

// HandlePatchThing todo
func HandlePatchThing(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {

}

// HandleGetObservedProperties todo
func HandleGetObservedProperties(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {

}

// HandleGetObservedPropertyById todo
func HandleGetObservedPropertyByID(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleGetObservedPropertyFromDatastream todo
func HandleGetObservedPropertyFromDatastream(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandlePostObservedProperty todo
func HandlePostObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleDeleteObservedProperty todo
func HandleDeleteObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandlePatchObservedProperty todo
func HandlePatchObservedProperty(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleGetLocations todo
func HandleGetLocations(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleGetLocationById todo
func HandleGetLocationByID(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandlePostLocation todo
func HandlePostLocation(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandlePostAndLinkLocation todo
func HandlePostAndLinkLocation(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleDeleteLocation todo
func HandleDeleteLocation(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandlePatchLocation todo
func HandlePatchLocation(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleGetDatastreams todo
func HandleGetDatastreams(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleGetDatastreamById todo
func HandleGetDatastreamByID(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleGetDatastreamsByThing todo
func HandleGetDatastreamsByThing(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandlePostDatastream todo
func HandlePostDatastream(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandlePostAndLinkDatastream todo
func HandlePostAndLinkDatastream(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleDeleteDatastream todo
func HandleDeleteDatastream(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandlePatchDatastream todo
func HandlePatchDatastream(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleGetSensors todo
func HandleGetSensors(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleGetSensorById todo
func HandleGetSensorByID(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandlePostSensors todo
func HandlePostSensors(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandleDeleteSensor todo
func HandleDeleteSensor(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
}

// HandlePatchSensor todo
func HandlePatchSensor(w http.ResponseWriter, r *http.Request, endpoint *Endpoint, api *API) {
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
func getQueryOptions(r *http.Request) (*QueryOptions, []error) {
	query := r.URL.Query()
	if len(query) == 0 {
		return nil, nil
	}

	qo, e := CreateQueryOptions(query)
	return qo, e
}

// tryParseEntity tries to parse a request body into the given entity
// calls SendError to handle the error if the body can't be parsed to the given entity
func tryParseEntity(w http.ResponseWriter, r io.ReadCloser, t interface{}) bool {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(t)
	if err != nil {
		sendError(w, http.StatusBadRequest, []error{err})
		return false
	}

	return true
}

// handleGetRequest is the default function to handle incoming GET requests
func handleGetRequest(w http.ResponseWriter, e *Endpoint, r *http.Request, h *func(q *QueryOptions) (interface{}, error), statusCode int) {
	queryOptions, err := getQueryOptions(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	_, errors := e.AreQueryOptionsSupported(queryOptions)
	if errors != nil {
		sendError(w, http.StatusBadRequest, errors)
		return
	}

	handler := *h
	data, err2 := handler(queryOptions)
	if err2 != nil {
		sendError(w, http.StatusInternalServerError, errors)
		return
	}

	sendJSONResponse(w, statusCode, data)
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
func sendError(w http.ResponseWriter, status int, error []error) {
	//errors cannot be marshalled, create strings
	errors := make([]string, len(error))
	for idx, value := range error {
		errors[idx] = fmt.Sprintf("%s", value)
	}

	statusText := http.StatusText(status)
	errorResponse := ErrorResponse{
		Error: ErrorContent{
			StatusText: statusText,
			StatusCode: status,
			Messages:   errors,
		},
	}

	sendJSONResponse(w, status, errorResponse)
}
