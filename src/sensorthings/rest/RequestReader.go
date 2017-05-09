package rest

import (
	"errors"
	"net/http"
	"strings"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/odata"
	"github.com/gorilla/mux"
	"strconv"
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
	//If request contains parameters from route wildcard convert it to a select query
	vars := mux.Vars(r)
	value := vars["params"]

	values := r.URL.Query()
	if len(vars["params"]) > 0 {
		//If $ref found create select query with id
		if vars["params"] == "$ref" {
			value = "id"
			values["$ref"] = []string{"true"}
		}

		values["$select"] = []string{value}
	}

	if strings.HasSuffix(r.URL.Path, "$value") {
		values["$value"] = []string{"true"}
	}

	if t, ok := r.URL.Query()["$top"]; !ok {
		values["$top"] = []string{strconv.Itoa(MaxEntities)}
	} else {
		top, err := strconv.Atoi(t[0])
		if err != nil || top > MaxEntities {
			values["$top"] = []string{strconv.Itoa(MaxEntities)}
		}
	}

	if _, ok := values["$skip"]; !ok {
		values["$skip"] = []string{"0"}
	}

	qo, e := odata.ParseUrlQuery(values)
	if e != nil {
		return nil, []error{e}
	}

	return qo, nil
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
