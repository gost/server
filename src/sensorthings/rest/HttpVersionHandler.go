package rest

import (
	"net/http"

	"github.com/geodan/gost/src/sensorthings/models"
)

// HandleVersion retrieves current version information and sends it back to the user
func HandleVersion(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	versionInfo := a.GetVersionInfo()
	sendJSONResponse(w, http.StatusOK, versionInfo, nil)
}
