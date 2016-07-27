package rest

import (
	"errors"
	"net/http"
	"strings"

	gostErrors "github.com/geodan/gost/src/errors"
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
