package handlers

import (
	"net/http"

	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
)

// HandlePostCreateObservations ...
func HandlePostCreateObservations(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	ob := &entities.CreateObservations{}
	handle := func() (interface{}, []error) { return a.PostCreateObservations(ob) }
	handlePostRequest(w, endpoint, r, ob, &handle, a.GetConfig().Server.IndentedJSON)
}
