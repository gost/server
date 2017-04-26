package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

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

			mVal := []byte{}
			for k, v := range m {
				if strings.ToLower(k) == qo.QuerySelect.Params[0] {
					mVal = v
				}
			}

			if len(mVal) == 0 {
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
	var b []byte
	var err error
	if IndentJSON {
		b, err = json.MarshalIndent(data, "", "   ")
	} else {
		b, err = json.Marshal(data)
	}

	// This code is needed if the response contains special characters like &, <, >, 
	// and those characters must not be converted to safe encoding.
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
