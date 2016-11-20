package rest

import (
	"errors"
	"net/http"
	"strings"

	"fmt"
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
	query := make(map[string]string)
	if strings.Contains(r.URL.String(), "?") {
		splitQuery := strings.Split(r.URL.RawQuery, "&")
		for _, sq := range splitQuery {
			splitIndex := strings.Index(sq, "=")
			if splitIndex == -1 {
				return nil, []error{fmt.Errorf("Incorrect request: %s", r.URL.RawQuery)}
			}

			query[sq[:splitIndex]] = sq[splitIndex+1:]
		}
	}

	//If request contains parameters from route wildcard convert it to a select query
	vars := mux.Vars(r)
	value := vars["params"]

	if len(vars["params"]) > 0 {
		//If $ref found create select query with id
		if vars["params"] == "$ref" {
			value = "id"
			query["$ref"] = "true"
		}

		query["$select"] = value
	}

	if strings.HasSuffix(r.URL.Path, "$value") {
		query["$value"] = "true"
	}

	// if $top is not found, retrieve max 200
	if _, ok := query["$top"]; !ok {
		query["$top"] = "200"
	}

	if _, ok := query["$skip"]; !ok {
		query["$skip"] = "0"
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
