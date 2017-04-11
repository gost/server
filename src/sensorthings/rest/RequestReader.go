package rest

import (
	"errors"
	"net/http"
	"strings"

	"fmt"
	"github.com/crestonbunch/godata"
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

	if len(vars["params"]) > 0 {
		//If $ref found create select query with id
		if vars["params"] == "$ref" {
			value = "id"
			r.URL.Query()["$ref"] = []string{"true"}
		}

		r.URL.Query()["$select"] = []string{value}
	}

	if strings.HasSuffix(r.URL.Path, "$value") {
		r.URL.Query()["$value"] = []string{"true"}
	}

	if t, ok := r.URL.Query()["$top"]; !ok {
		r.URL.Query()["$top"] = []string{strconv.Itoa(MaxEntities)}
	} else {
		top, err := strconv.Atoi(t[0])
		if err != nil || top > MaxEntities {
			r.URL.Query()["$top"] = []string{strconv.Itoa(MaxEntities)}
		}
	}

	if _, ok := r.URL.Query()["$skip"]; !ok {
		r.URL.Query()["$skip"] = []string{"0"}
	}

	qo, e := godata.ParseUrlQuery(r.URL.Query())
	vals := *qo.Filter.Tree
	fmt.Printf("%v", vals)

	return nil, []error{e}
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
