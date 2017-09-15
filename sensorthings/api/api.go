package api

import (
	"fmt"
	"strings"

	entities "github.com/gost/core"
	"github.com/gost/server/configuration"
	gostErrors "github.com/gost/server/errors"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/mqtt"
	"github.com/gost/server/sensorthings/odata"
	"github.com/gost/server/sensorthings/rest/config"
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

	odata.SupportedExpandParameters = expandParams
	odata.SupportedSelectParameters = selectParams
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
func (a *APIv1) GetBasePathInfo() *entities.ArrayResponse {
	bpi := []models.Endpoint{}
	ep := *a.GetEndpoints()
	for _, e := range ep {
		if e.ShowOutputInfo() {
			bpi = append(bpi, e)
		}
	}

	var i interface{} = bpi
	basePathInfo := entities.ArrayResponse{
		Data: &i,
	}

	return &basePathInfo
}

// GetEndpoints returns all configured endpoints for the HTTP server
func (a *APIv1) GetEndpoints() *map[entities.EntityType]models.Endpoint {
	if a.endPoints == nil {
		a.endPoints = config.CreateEndPoints(a.config.GetExternalServerURI())
	}

	return &a.endPoints
}

func (a *APIv1) initRest() {
	if a.config.Server.MaxEntityResponse == 0 {
		a.config.Server.MaxEntityResponse = configuration.DefaultMaxEntries
	}
}

// GetTopics returns all configured topics for the MQTT client
func (a *APIv1) GetTopics() *[]models.Topic {
	if a.topics == nil {
		a.topics = mqtt.CreateTopics()
	}

	return &a.topics
}

// SetLinks processes the entities by setting the necessary links before sending back
func (a *APIv1) SetLinks(entity entities.Entity, qo *odata.QueryOptions) {
	// a $ref request, id's are selected to create selfLink, remove after setting self url
	if qo != nil && qo.Ref != nil && bool(*qo.Ref) {
		entity.SetSelfLink(a.config.GetExternalServerURI())
		entity.SetID(nil)
	} else if qo == nil || qo.Select == nil || len(qo.Select.SelectItems) == 0 { //no query options, set all links
		entity.SetAllLinks(a.config.GetExternalServerURI())
	}
}

// CreateNextLink creates the link to the next page with results
//  incomingUrl is the url of the request excluding oData query params
func (a *APIv1) CreateNextLink(incomingURL string, qo *odata.QueryOptions) string {
	// do not create a nextLink when there is no top and skip given
	if qo == nil || qo.Top == nil || qo.Skip == nil || (int(*qo.Top) == 0 && int(*qo.Skip) == 0) {
		// todo: check unreachable code here?
		return ""
	}

	queryString := ""
	if qo.Filter != nil {
		queryString = appendQueryPart(queryString, fmt.Sprintf("$filter=%s", qo.RawFilter))
	}
	if qo.Count != nil {
		queryString = appendQueryPart(queryString, fmt.Sprintf("$count=%v", *qo.Count))
	}
	if qo.Expand != nil {
		queryString = appendQueryPart(queryString, fmt.Sprintf("$expand=%s", qo.RawExpand))
	}
	if qo.OrderBy != nil {
		queryString = appendQueryPart(queryString, fmt.Sprintf("$orderby=%s", qo.RawOrderBy))
	}
	if qo.Format != nil {
		queryString = appendQueryPart(queryString, fmt.Sprintf("$format=%s", qo.Format))
	}
	if qo.Top != nil {
		queryString = appendQueryPart(queryString, fmt.Sprintf("$top=%v", *qo.Top))
	}
	if qo.Skip != nil {
		queryString = appendQueryPart(queryString, fmt.Sprintf("$skip=%v", int(*qo.Skip)+int(*qo.Top)))
	}

	return fmt.Sprintf("%s%s", incomingURL, queryString)
}

func containsMandatoryParams(entity interface{}) (bool, []error) {
	contains := false
	var errors []error

	if entity != nil {
		switch e := entity.(type) {
		case *entities.Thing:
			contains, errors = e.ContainsMandatoryParams()
		case *entities.Location:
			contains, errors = e.ContainsMandatoryParams()
		case *entities.HistoricalLocation:
			contains, errors = e.ContainsMandatoryParams()
		case *entities.Datastream:
			contains, errors = e.ContainsMandatoryParams()
		case *entities.Sensor:
			contains, errors = e.ContainsMandatoryParams()
		case *entities.ObservedProperty:
			contains, errors = e.ContainsMandatoryParams()
		case *entities.Observation:
			contains, errors = e.ContainsMandatoryParams()
		case *entities.FeatureOfInterest:
			contains, errors = e.ContainsMandatoryParams()
		}
	}

	// Wrap errors in BadRequest
	if errors != nil {
		for i := range errors {
			errors[i] = gostErrors.NewBadRequestError(errors[i])
		}
	}

	return contains, errors
}

// createArrayResponse creates the ArrayResponse to send back to the user
func (a *APIv1) createArrayResponse(count int, hasNext bool, path string, qo *odata.QueryOptions, data interface{}) *entities.ArrayResponse {
	ar := &entities.ArrayResponse{
		Data: &data,
	}

	if hasNext {
		ar.NextLink = a.CreateNextLink(path, qo)
	}

	if qo != nil && qo.Count != nil && bool(*qo.Count) == true {
		ar.Count = count
	}

	return ar
}

func appendQueryPart(base string, q string) string {
	prefix := "?"
	if strings.Contains(base, "?") {
		prefix = "&"
	}

	return fmt.Sprintf("%s%s%s", base, prefix, q)
}
