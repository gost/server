package rest

import (
	"net/http"

	"github.com/geodan/gost/sensorthings/models"
)

// HandleAPIRoot will return a JSON array of the available SensorThings resource endpoints.
func HandleAPIRoot(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	bpi := a.GetBasePathInfo()
	sendJSONResponse(w, http.StatusOK, bpi, nil)
}
