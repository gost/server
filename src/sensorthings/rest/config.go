package rest

import "github.com/geodan/gost/src/sensorthings/models"

// CreateEndPoints creates the pre-defined endpoint config, the config contains all endpoint info
// describing the SupportedQueryOptions (if needed) and EndpointOperation for each endpoint
// parameter externalURL is the URL where the GOST service can be reached, main endpoint urls
// are generated based upon this URL
func CreateEndPoints(externalURL string) []models.Endpoint {
	endpoints := []models.Endpoint{
		createVersionEndpoint(externalURL),
		createRootEndpoint(externalURL),
		createThingsEndpoint(externalURL),
		createDatastreamsEndpoint(externalURL),
		createObservedPropertiesEndpoint(externalURL),
		createLocationsEndpoint(externalURL),
		createSensorsEndpoint(externalURL),
		createObservationsEndpoint(externalURL),
		createFeaturesOfInterestEndpoint(externalURL),
		createHistoricalLocationsEndpoint(externalURL),
	}

	return endpoints
}
