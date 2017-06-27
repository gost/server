package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
)

// handlePatchRequest todo: currently almost same as handlePostRequest, merge if it stays like this
func handlePatchRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, entity entities.Entity, h *func() (interface{}, error)) {
	if !checkContentType(w, r) {
		return
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	byteData, _ := ioutil.ReadAll(r.Body)
	err := entity.ParseEntity(byteData)
	if err != nil {
		sendError(w, []error{err})
		return
	}

	handle := *h
	data, err2 := handle()
	if err2 != nil {
		sendError(w, []error{err2})
		return
	}

	w.Header().Add("Location", entity.GetSelfLink())

	sendJSONResponse(w, http.StatusOK, data, nil)
}
