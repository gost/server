package handlers

import (
	"net/http"

	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/rest/writer"
)

// handleDeleteRequest
func handleDeleteRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, h *func() error, indentJSON bool) {
	handle := *h
	err := handle()
	if err != nil {
		writer.SendError(w, []error{err}, indentJSON)
		return
	}

	writer.SendJSONResponse(w, http.StatusOK, nil, nil, indentJSON)
}
