package handlers

import (
	"net/http"

	"fmt"

	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/odata"
	"github.com/gost/server/sensorthings/rest/writer"
)

// handleGetRequest is the default function to handle incoming GET requests
func handleGetRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, h *func(q *odata.QueryOptions, path string) (interface{}, error), indentJSON bool, maxEntities int, externalURI string) {
	// Parse query options from request
	queryOptions, err := odata.GetQueryOptions(r, maxEntities)
	if err != nil && len(err) > 0 {
		writer.SendError(w, err, indentJSON)
		return
	}

	// Run the handler func such as Api.GetThingById
	handler := *h
	data, err2 := handler(queryOptions, fmt.Sprintf(externalURI+r.URL.RawPath))
	if err2 != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err2.Error()))
		return
	}

	writer.SendJSONResponse(w, http.StatusOK, data, queryOptions, indentJSON)
}
