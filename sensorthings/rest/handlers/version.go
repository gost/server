package handlers

import (
	"net/http"

	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/rest/writer"
)

// HandleVersion retrieves current version information and sends it back to the user
func HandleVersion(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
	a := *api
	versionInfo := a.GetVersionInfo()
	writer.SendJSONResponse(w, http.StatusOK, versionInfo, nil, a.GetConfig().Server.IndentedJSON)
}
