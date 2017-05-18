package rest

import (
	"fmt"

	"github.com/geodan/gost/sensorthings/models"
)

func createVersionEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Version",
		OutputInfo: false,
		URL:        fmt.Sprintf("%s/%s", externalURL, "Version"),
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/version", HandleVersion},
		},
	}
}
