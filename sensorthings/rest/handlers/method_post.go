package handlers

import (
	"net/http"

	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/reader"
	"github.com/gost/server/sensorthings/rest/writer"
)

// handlePostRequest
func handlePostRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, entity entities.Entity, h *func() (interface{}, []error), indentJSON bool) {
	if !reader.CheckContentType(w, r, indentJSON) {
		return
	}

	byteData := reader.CheckAndGetBody(w, r, indentJSON)
	if byteData == nil {
		return
	}

	err := reader.ParseEntity(entity, byteData)
	if err != nil {
		writer.SendError(w, []error{err}, indentJSON)
		return
	}

	handle := *h
	data, err2 := handle()
	if err2 != nil {
		writer.SendError(w, err2, indentJSON)
		return
	}

	w.Header().Add("Location", entity.GetSelfLink())

	writer.SendJSONResponse(w, http.StatusCreated, data, nil, indentJSON)
}
