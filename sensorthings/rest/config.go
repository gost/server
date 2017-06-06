package rest

import (
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
)

// Endpoints contains the current set-up of the endpoints
var Endpoints = map[entities.EntityType]models.Endpoint{}

// MaxEntities is used to determine the max entities to return to the client
var MaxEntities = 200

// IndentJSON is used to determine if the JSON response should be indented or not
var IndentJSON = true

// ExternalURI is used for defining an external url which is shown to the requester instead of the internal uri
var ExternalURI = ""

// CreateEndPoints creates the pre-defined endpoint config, the config contains all endpoint info
// describing the SupportedQueryOptions (if needed) and EndpointOperation for each endpoint
// parameter externalURL is the URL where the GOST service can be reached, main endpoint urls
// are generated based upon this URL
func CreateEndPoints(externalURL string) map[entities.EntityType]models.Endpoint {
	Endpoints = map[entities.EntityType]models.Endpoint{
		entities.EntityTypeVersion:            createVersionEndpoint(externalURL),
		entities.EntityTypeUnknown:            createRootEndpoint(externalURL),
		entities.EntityTypeThing:              createThingsEndpoint(externalURL),
		entities.EntityTypeDatastream:         createDatastreamsEndpoint(externalURL),
		entities.EntityTypeObservedProperty:   createObservedPropertiesEndpoint(externalURL),
		entities.EntityTypeLocation:           createLocationsEndpoint(externalURL),
		entities.EntityTypeSensor:             createSensorsEndpoint(externalURL),
		entities.EntityTypeObservation:        createObservationsEndpoint(externalURL),
		entities.EntityTypeFeatureOfInterest:  createFeaturesOfInterestEndpoint(externalURL),
		entities.EntityTypeHistoricalLocation: createHistoricalLocationsEndpoint(externalURL),
	}

	return Endpoints
}
