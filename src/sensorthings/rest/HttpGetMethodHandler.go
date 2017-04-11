package rest

import (
	"net/http"

	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// handleGetRequest is the default function to handle incoming GET requests
func handleGetRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, h *func(q *odata.QueryOptions, path string) (interface{}, error)) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	// Parse query options from request
	queryOptions, err := getQueryOptions(r)
	if err != nil {
		sendError(w, err)
		return
	}

	// Check if the requested endpoints supports the parsed queries
	endpoint := *e
	_, err = endpoint.AreQueryOptionsSupported(queryOptions)
	if err != nil {
		sendError(w, err)
		return
	}

	// Run the handler func such as Api.GetThingById
	handler := *h
	data, err2 := handler(queryOptions, fmt.Sprintf(ExternalURI+r.URL.RawPath))
	if err2 != nil {
		sendError(w, []error{err2})
		return
	}

	sendJSONResponse(w, http.StatusOK, data, queryOptions)
}
