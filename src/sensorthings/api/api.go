package api

import (
	"github.com/geodan/gost/src/configuration"

	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/mqtt"
	"github.com/geodan/gost/src/sensorthings/rest"
)

// APIv1 is the default implementation of SensorThingsApi, API needs a database
// provider, config, endpoint information to setup te needed services
type APIv1 struct {
	db        models.Database
	config    configuration.Config
	endPoints []models.Endpoint
	topics    []models.Topic
	mqtt      models.MQTTClient
}

// NewAPI Initialise a new SensorThings API
func NewAPI(database models.Database, config configuration.Config, mqtt models.MQTTClient) models.API {
	return &APIv1{
		db:     database,
		mqtt:   mqtt,
		config: config,
	}
}

// GetConfig return the current configuration.Config set for the api
func (a *APIv1) GetConfig() *configuration.Config {
	return &a.config
}

// GetVersionInfo retrieves the version info of the current supported SensorThings API Version and running server version
func (a *APIv1) GetVersionInfo() *models.VersionInfo {
	versionInfo := models.VersionInfo{
		GostServerVersion: models.GostServerVersion{Version: configuration.ServerVersion},
		APIVersion:        models.APIVersion{Version: configuration.SensorThingsAPIVersion},
	}

	return &versionInfo
}

// GetBasePathInfo when navigating to the base resource path will return a JSON array of the available SensorThings resource endpoints.
func (a *APIv1) GetBasePathInfo() *models.ArrayResponse {
	bpi := []models.Endpoint{}
	ep := *a.GetEndpoints()
	for _, e := range ep {
		if e.ShowOutputInfo() {
			bpi = append(bpi, e)
		}
	}

	var i interface{} = bpi
	basePathInfo := models.ArrayResponse{
		Data: &i,
	}

	return &basePathInfo
}

// GetEndpoints returns all configured endpoints for the HTTP server
func (a *APIv1) GetEndpoints() *[]models.Endpoint {
	if a.endPoints == nil {
		a.endPoints = rest.CreateEndPoints(a.config.GetExternalServerURI())
	}

	return &a.endPoints
}

// GetTopics returns all configured topics for the MQTT client
func (a *APIv1) GetTopics() *[]models.Topic {
	if a.topics == nil {
		a.topics = mqtt.CreateTopics()
	}

	return &a.topics
}
