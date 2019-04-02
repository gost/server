package config

import (
	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
)

// Endpoints contains the current set-up of the endpoints
var Endpoints = map[entities.EntityType]models.Endpoint{}

// CreateEndPoints creates the pre-defined endpoint config, the config contains all endpoint info
// describing the SupportedQueryOptions (if needed) and EndpointOperation for each endpoint
// parameter externalURL is the URL where the GOST service can be reached, main endpoint urls
// are generated based upon this URL
func CreateEndPoints(externalURL string) map[entities.EntityType]models.Endpoint {
	Endpoints = map[entities.EntityType]models.Endpoint{
		entities.EntityTypeVersion:            CreateVersionEndpoint(externalURL),
		entities.EntityTypeUnknown:            CreateRootEndpoint(externalURL),
		entities.EntityTypeThing:              CreateThingsEndpoint(externalURL),
		entities.EntityTypeDatastream:         CreateDatastreamsEndpoint(externalURL),
		entities.EntityTypeObservedProperty:   CreateObservedPropertiesEndpoint(externalURL),
		entities.EntityTypeLocation:           CreateLocationsEndpoint(externalURL),
		entities.EntityTypeSensor:             CreateSensorsEndpoint(externalURL),
		entities.EntityTypeObservation:        CreateObservationsEndpoint(externalURL),
		entities.EntityTypeFeatureOfInterest:  CreateFeaturesOfInterestEndpoint(externalURL),
		entities.EntityTypeHistoricalLocation: CreateHistoricalLocationsEndpoint(externalURL),
		entities.EntityTypeCreateObservations: CreateCreateObservationsEndpoint(externalURL),
	}

	return Endpoints
}
