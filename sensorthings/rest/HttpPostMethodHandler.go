package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
)

// handlePostRequest
func handlePostRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, entity entities.Entity, h *func() (interface{}, []error)) {
	if !checkContentType(w, r) {
		return
	}

	byteData, _ := ioutil.ReadAll(r.Body)
	err := entity.ParseEntity(byteData)
	if err != nil {
		sendError(w, []error{err})
		return
	}

	handle := *h
	data, err2 := handle()
	if err2 != nil {
		sendError(w, err2)
		return
	}

	w.Header().Add("Location", entity.GetSelfLink())

	sendJSONResponse(w, http.StatusCreated, data, nil)
}
