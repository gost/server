package handlers

import (
	"net/http"

	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/rest/writer"
)

// HandleAPIRoot will return a JSON array of the available SensorThings resource endpoints.
func HandleAPIRoot(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	bpi := a.GetBasePathInfo()
	writer.SendJSONResponse(w, http.StatusOK, bpi, nil, a.GetConfig().Server.IndentedJSON)
}
