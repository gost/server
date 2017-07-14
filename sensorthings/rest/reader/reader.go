package reader

import (
	"errors"
	"net/http"
	"strings"

	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/rest/writer"
	"github.com/gorilla/mux"
)

// GetEntityID retrieves the id from the request, for example
// the request http://mysensor.com/V1.0/Things(1236532) returns 1236532 as id
func GetEntityID(r *http.Request) string {
	vars := mux.Vars(r)
	value := vars["id"]
	substring := value[1 : len(value)-1]
	return substring
}

func CheckContentType(w http.ResponseWriter, r *http.Request, indentJSON bool) bool {
	// maybe needs to add case-insentive check?
	if len(r.Header.Get("Content-Type")) > 0 {
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			writer.SendError(w, []error{gostErrors.NewBadRequestError(errors.New("Missing or wrong Content-Type, accepting: application/json"))}, indentJSON)
			return false
		}
	}

	return true
}
