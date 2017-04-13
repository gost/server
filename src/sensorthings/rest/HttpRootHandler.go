package rest

import (
	"net/http"

	"github.com/geodan/gost/src/sensorthings/models"
)

// HandleAPIRoot will return a JSON array of the available SensorThings resource endpoints.
func HandleAPIRoot(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	bpi := a.GetBasePathInfo()
	w.Header().Add("Access-Control-Allow-Origin", "*")
	sendJSONResponse(w, http.StatusOK, bpi, nil)
}
