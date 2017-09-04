package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"errors"
	gostErrors "github.com/gost/server/errors"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/odata"
)

// SendJSONResponse sends the desired message to the user
// the message will be marshalled into JSON
func SendJSONResponse(w http.ResponseWriter, status int, data interface{}, qo *odata.QueryOptions, indentJSON bool) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if data != nil {
		b, err := JSONMarshal(data, true, indentJSON)
		if err != nil {
			panic(err)
		}

		if qo != nil {
			// $count for collection is requested url/$count and not the query ?$count=true, ToDo: move to API code?e
			if qo.CollectionCount != nil && bool(*qo.CollectionCount) == true {
				var m map[string]json.RawMessage
				json.Unmarshal(b, &m)
				if count, ok := m["@iot.count"]; ok {
					b = []byte(string(count))
				} else {
					SendError(w, []error{gostErrors.NewBadRequestError(errors.New("/$count not available for endpoint"))}, indentJSON)
					return
				}
			} else if qo.Value != nil && bool(*qo.Value) == true {
				// $value is requested only send back the value
				errMessage := fmt.Errorf("Unable to retrieve $value for %v", qo.Select.SelectItems)
				var m map[string]json.RawMessage
				err = json.Unmarshal(b, &m)
				if err != nil || qo.Select == nil || qo.Select.SelectItems == nil || len(qo.Select.SelectItems) == 0 {
					SendError(w, []error{gostErrors.NewRequestInternalServerError(errMessage)}, indentJSON)
					return
				}

				// if selected equals the key in json add to mVal
				mVal := []byte{}
				for k, v := range m {
					if strings.ToLower(k) == qo.Select.SelectItems[0].Segments[0].Value {
						mVal = v
					}
				}

				if len(mVal) == 0 {
					SendError(w, []error{gostErrors.NewBadRequestError(errMessage)}, indentJSON)
					return
				}

				value := string(mVal[:])
				value = strings.TrimPrefix(value, "\"")
				value = strings.TrimSuffix(value, "\"")

				b = []byte(value)
			}
		}

		w.WriteHeader(status)
		w.Write(b)
	}
}

//JSONMarshal converts the data and converts special characters such as &
func JSONMarshal(data interface{}, safeEncoding, indentJSON bool) ([]byte, error) {
	var b []byte
	var err error
	if indentJSON {
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

// SendError creates an ErrorResponse message and sets it to the user
// using SendJSONResponse
func SendError(w http.ResponseWriter, error []error, indentJSON bool) {
	//errors cannot be marshalled, create strings
	errors := make([]string, len(error))
	for idx, value := range error {
		errors[idx] = value.Error()
	}

	// Set the status code, default 500 for error, check if there is an ApiError an get
	var statusCode = http.StatusInternalServerError

	if len(error) > 0 {
		// if there is Encoding type error, sends bad request (400 range)
		if strings.Contains(errors[0], "Encoding not supported") {
			statusCode = http.StatusBadRequest
		}

		switch e := error[0].(type) {
		case gostErrors.APIError:
			statusCode = e.GetHTTPErrorStatusCode()
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

	SendJSONResponse(w, statusCode, errorResponse, nil, indentJSON)
}
