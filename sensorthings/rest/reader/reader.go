package reader

import (
	"errors"
	"net/http"
	"strings"

	"fmt"
	"io/ioutil"

	"github.com/gorilla/mux"
	gostErrors "github.com/gost/server/errors"
	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/rest/writer"
)

// GetEntityID retrieves the id from the request, for example
// the request http://mysensor.com/V1.0/Things(1236532) returns 1236532 as id
func GetEntityID(r *http.Request) string {
	vars := mux.Vars(r)
	value := vars["id"]
	substring := value[1 : len(value)-1]
	return substring
}

// CheckContentType checks if there is a content-type header, if so check if it is of type
// application/json, if not return an error, SensorThings server only accepts application/json
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

// CheckAndGetBody checks if the request body is not nil and tries to read it in a byte slice
// when an error occurs an error will be send back using the ResponseWriter
func CheckAndGetBody(w http.ResponseWriter, r *http.Request, indentJSON bool) []byte {
	if r.Body == nil {
		writer.SendError(w, []error{gostErrors.NewBadRequestError(fmt.Errorf("No body found in request"))}, indentJSON)
		return nil
	}

	byteData, _ := ioutil.ReadAll(r.Body)
	return byteData
}

// ParseEntity tries to convert the byte data into the given interface of type entity
// if an error returns it will be wraped inside an gosterror
func ParseEntity(entity entities.Entity, data []byte) error {
	var err error

	err = entity.ParseEntity(data)

	if err != nil {
		err = gostErrors.NewBadRequestError(err)
	}

	return err
}
