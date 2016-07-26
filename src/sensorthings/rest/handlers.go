package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"bytes"
	"fmt"
	"strings"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
	"github.com/gorilla/mux"
)


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

func checkContentType(w http.ResponseWriter, r *http.Request) bool {
	// maybe needs to add case-insentive check?
	if len(r.Header.Get("Content-Type")) > 0 {
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			sendError(w, []error{gostErrors.NewBadRequestError(errors.New("Missing or wrong Content-Type, accepting: application/json"))})
			return false
		}
	}

	return true
}

// sendJSONResponse sends the desired message to the user
// the message will be marshalled into an indented JSON format
func sendJSONResponse(w http.ResponseWriter, status int, data interface{}, qo *odata.QueryOptions) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if data != nil {
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
