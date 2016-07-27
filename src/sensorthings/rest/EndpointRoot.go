package rest

import (
	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
)

func createRootEndpoint(externalURL string) *Endpoint {
	return &Endpoint{
		Name:       "Root",
		OutputInfo: false,
		URL:        fmt.Sprintf("%s/%s", externalURL, "v1.0"),
		Operations: []models.EndpointOperation{
			{models.HTTPOperationGet, "/v1.0", HandleAPIRoot},
		},
	}
}
