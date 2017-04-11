package api

import (
	"fmt"
	"strings"

	"github.com/geodan/gost/src/configuration"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/mqtt"
	"github.com/geodan/gost/src/sensorthings/odata"
	"github.com/geodan/gost/src/sensorthings/rest"
)

// APIv1 is the default implementation of SensorThingsApi, API needs a database
// provider, config, endpoint information to setup te needed services
type APIv1 struct {
	db            models.Database
	config        configuration.Config
	endPoints     map[entities.EntityType]models.Endpoint
	topics        []models.Topic
	mqtt          models.MQTTClient
	acceptedPaths []string
}

// NewAPI Initialise a new SensorThings API
func NewAPI(database models.Database, config configuration.Config, mqtt models.MQTTClient) models.API {
	api := &APIv1{
		db:     database,
		mqtt:   mqtt,
		config: config,
		acceptedPaths: []string{
			"v1.0",
			"thing",
			"things",
			"datastream",
			"datastreams",
			"location",
			"locations",
			"historicallocation",
			"historicallocations",
			"sensor",
			"sensors",
			"observation",
			"observations",
			"observedproperty",
			"observedproperties",
			"featureofinterest",
			"featurseofinterest",
			"$value",
			"dashboard",
		},
	}
	api.initRest()
	api.Start()
	return api
}

// Start is used to set the initial state of the api such as loading of the foi states
func (a *APIv1) Start() {
	eps := *a.GetEndpoints()
	expandParams := map[string][]string{}
	selectParams := map[string][]string{}
	for _, e := range eps {
		expandParams[e.GetName()] = e.GetSupportedExpandParams()
		selectParams[e.GetName()] = e.GetSupportedSelectParams()
	}

	odata.Init(expandParams, selectParams)
}

// GetConfig return the current configuration.Config set for the api
func (a *APIv1) GetConfig() *configuration.Config {
	return &a.config
}

// GetAcceptedPaths returns an array of accepted endpoint paths
func (a *APIv1) GetAcceptedPaths() []string {
	return a.acceptedPaths
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
func (a *APIv1) GetEndpoints() *map[entities.EntityType]models.Endpoint {
	if a.endPoints == nil {
		a.endPoints = rest.CreateEndPoints(a.config.GetExternalServerURI())
	}

	return &a.endPoints
}

func (a *APIv1) initRest() {
	if a.config.Server.MaxEntityResponse == 0 {
		a.config.Server.MaxEntityResponse = configuration.DefaultMaxEntries
	}

	rest.ExternalURI = a.config.GetExternalServerURI()
	rest.IndentJSON = a.config.Server.IndentedJSON
	rest.MaxEntities = a.config.Server.MaxEntityResponse
}

// GetTopics returns all configured topics for the MQTT client
func (a *APIv1) GetTopics() *[]models.Topic {
	if a.topics == nil {
		a.topics = mqtt.CreateTopics()
	}

	return &a.topics
}

// QueryOptionsSupported checks if the query options are supported for the current entity
func (a *APIv1) QueryOptionsSupported(qo *odata.QueryOptions, entity entities.Entity) (bool, error) {
	if qo == nil {
		return true, nil
	}

	return true, nil
}

// ProcessGetRequest processes the entities by setting the necessary links before sending back
func (a *APIv1) ProcessGetRequest(entity entities.Entity, qo *odata.QueryOptions) {
	// a $ref request, id's are selected to create selfLink, remove after setting self url
	if qo != nil && qo.QueryOptionRef {
		entity.SetSelfLink(a.config.GetExternalServerURI())
		entity.SetID(nil)
	} else if qo == nil || qo.QuerySelect.IsNil() || len(qo.QuerySelect.Params) == 0 { //no query options, set all links
		entity.SetAllLinks(a.config.GetExternalServerURI())
	}
}

// CreateNextLink creates the link to the next page with results
//  count is the number of total entities in the current query
//  incomingUrl is the url of the request excluding oData query params
func (a *APIv1) CreateNextLink(count int, incomingURL string, qo *odata.QueryOptions) string {
	// do not create a nextLink when there is no top and skip given
	if qo == nil && qo.QueryTop.Limit == 0 && qo.QuerySkip.Index == 0 {
		return ""
	}

	// do not create a nextLink when the current page has no following one
	if qo.QueryTop.Limit+qo.QuerySkip.Index >= count || count < qo.QueryTop.Limit {
		return ""
	}

	queryString := ""
	if !qo.QueryFilter.IsNil() {
		queryString = appendQueryPart(queryString, fmt.Sprintf("%s=%s", odata.QueryOptionFilter.String(), qo.QueryFilter.RawQuery))
	}
	if !qo.QueryCount.IsNil() {
		queryString = appendQueryPart(queryString, fmt.Sprintf("%s=%s", odata.QueryOptionCount.String(), qo.QueryCount.RawQuery))
	}
	if !qo.QueryExpand.IsNil() {
		queryString = appendQueryPart(queryString, fmt.Sprintf("%s=%s", odata.QueryOptionExpand.String(), qo.QueryExpand.RawQuery))
	}
	if !qo.QueryOrderBy.IsNil() {
		queryString = appendQueryPart(queryString, fmt.Sprintf("%s=%s", odata.QueryOptionOrderBy.String(), qo.QueryOrderBy.RawQuery))
	}
	if !qo.QueryResultFormat.IsNil() {
		queryString = appendQueryPart(queryString, fmt.Sprintf("%s=%s", odata.QueryOptionResultFormat.String(), qo.QueryResultFormat.RawQuery))
	}
	if !qo.QueryTop.IsNil() {
		queryString = appendQueryPart(queryString, fmt.Sprintf("%s=%s", odata.QueryOptionTop.String(), qo.QueryTop.RawQuery))
	}
	if !qo.QuerySkip.IsNil() {
		queryString = appendQueryPart(queryString, fmt.Sprintf("%s=%v", odata.QueryOptionSkip.String(), qo.QuerySkip.Index+qo.QueryTop.Limit))
	}

	return fmt.Sprintf("%s/%s", incomingURL, queryString)
}

func appendQueryPart(base string, q string) string {
	prefix := "?"
	if strings.Contains(base, "?") {
		prefix = "&"
	}

	return fmt.Sprintf("%s%s%s", base, prefix, q)
}
