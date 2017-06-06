package rest

import (
	"net/http"

	"fmt"

	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// handleGetRequest is the default function to handle incoming GET requests
func handleGetRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, h *func(q *odata.QueryOptions, path string) (interface{}, error)) {

	// Parse query options from request
	queryOptions, err := getQueryOptions(r)
	if err != nil && len(err) > 0 {
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
