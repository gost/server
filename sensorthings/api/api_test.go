package api

import (
	"github.com/gost/server/configuration"
	"github.com/gost/server/database/postgis"
	gostErrors "github.com/gost/server/errors"
	"github.com/gost/server/mqtt"
	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/odata"

	"fmt"
	"strings"
	"testing"

	"github.com/gost/godata"
	"github.com/stretchr/testify/assert"
)

func TestCreateApi(t *testing.T) {
	// arrange
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100, 200)
	stAPI := NewAPI(database, cfg, mqttServer)

	// act
	ep := stAPI.GetEndpoints()
	config := stAPI.GetConfig()
	paths := stAPI.GetAcceptedPaths()
	versionInfo := stAPI.GetVersionInfo()
	basePathInfo := stAPI.GetBasePathInfo()

	endpoints := *ep

	// assert
	assert.NotNil(t, basePathInfo)
	assert.NotNil(t, versionInfo)
	assert.NotNil(t, paths)
	assert.NotNil(t, config)
	assert.NotNil(t, ep)
	assert.NotEqual(t, len(endpoints), 0, "Endpoints empty")
}

func TestGetTopics(t *testing.T) {
	// arrange
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100, 200)
	stAPI := NewAPI(database, cfg, mqttServer)

	// act
	topics := stAPI.GetTopics()

	// assert
	assert.NotNil(t, topics)
	firsttopic := (*topics)[0]
	assert.True(t, firsttopic.Path == "GOST/#")
}

func TestAppendQueryPart(t *testing.T) {
	// act
	result := appendQueryPart("base", "q")
	result1 := appendQueryPart("base?", "q")

	// assert
	assert.True(t, result == "base?q")
	assert.True(t, result1 == "base?&q")

}

func TestSetLinks(t *testing.T) {
	// arrange
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100, 200)
	stAPI := NewAPI(database, cfg, mqttServer)
	ds := entities.Datastream{}

	// act
	stAPI.SetLinks(&ds, nil)

	// assert
	assert.True(t, ds.GetSelfLink() == "/v1.0/Datastreams")
}

func TestSetLinkWithQuery(t *testing.T) {
	// arrange
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100, 200)
	stAPI := NewAPI(database, cfg, mqttServer)
	ds := entities.Datastream{}

	qo := &odata.QueryOptions{}
	qt, _ := godata.ParseTopString("2")
	qo.Top = qt
	ref := odata.GoDataRefQuery(true)
	qo.Ref = &ref
	// act
	stAPI.SetLinks(&ds, qo)

	// assert
	assert.True(t, ds.GetSelfLink() == "/v1.0/Datastreams")
	assert.True(t, ds.ID == nil)

}

func TestCreateNextLink(t *testing.T) {
	// arrange
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100, 200)
	stAPI := NewAPI(database, cfg, mqttServer)
	qo := &odata.QueryOptions{}

	qt, _ := godata.ParseTopString("2")
	qo.Top = qt
	qs, _ := godata.ParseSkipString("1")
	qo.Skip = qs

	// act
	result := stAPI.CreateNextLink(1, "http://www.nu.nl", qo)
	assert.NotNil(t, result)
	assert.True(t, result == "")

	qt, _ = godata.ParseTopString("1")
	qo.Top = qt
	filterString := "id eq 1"
	qf, err := godata.ParseFilterString(filterString)
	if err != nil {
		t.Errorf("Error parsing filter string: %v", err)
	}

	qo.Filter = qf
	qo.RawFilter = filterString
	// add QueryCount, QueryExpand, QueryOrderBy, QueryResultFormat

	result1 := stAPI.CreateNextLink(10, "http://www.nu.nl", qo)
	t.Logf("%v", result1)
	// assert
	assert.NotNil(t, result1)
	assert.True(t, strings.Contains(result1, "id eq 1"))
}

func TestCreateArrayResponseWithCount(t *testing.T) {
	// arrange
	testAPI := &APIv1{}
	count := 1
	path := "testPath"
	data := "test"
	countQuery := godata.GoDataCountQuery(true)
	qo := &odata.QueryOptions{}
	qo.Count = &countQuery

	// act
	arrayResponse := testAPI.createArrayResponse(count, path, qo, data)

	// assert
	assert.Equal(t, count, arrayResponse.Count)
	assert.Equal(t, data, fmt.Sprintf("%v", *arrayResponse.Data))
	assert.Equal(t, testAPI.CreateNextLink(count, path, qo), arrayResponse.NextLink)
}

func TestCreateArrayResponseWithoutCount(t *testing.T) {
	// arrange
	testAPI := &APIv1{}
	count := 10
	qo := &odata.QueryOptions{}

	// act
	arrayResponse := testAPI.createArrayResponse(count, "", qo, "")

	// assert
	assert.Equal(t, 0, arrayResponse.Count)
}

func TestContainsMandatoryParamsReturningBadRequest(t *testing.T) {
	// arrange
	thing := &entities.Thing{}
	location := &entities.Location{}
	historicalLocation := &entities.HistoricalLocation{}
	datastream := &entities.Datastream{}
	sensor := &entities.Sensor{}
	observedProperty := &entities.ObservedProperty{}
	observation := &entities.Observation{}
	featureOfinterest := &entities.FeatureOfInterest{}

	// act
	_, tErr := containsMandatoryParams(thing)
	_, lErr := containsMandatoryParams(location)
	_, hlErr := containsMandatoryParams(historicalLocation)
	_, dErr := containsMandatoryParams(datastream)
	_, sErr := containsMandatoryParams(sensor)
	_, opErr := containsMandatoryParams(observedProperty)
	_, oErr := containsMandatoryParams(observation)
	_, fErr := containsMandatoryParams(featureOfinterest)

	// assert
	assert.Equal(t, 400, getStatusCode(tErr))
	assert.Equal(t, 400, getStatusCode(lErr))
	assert.Equal(t, 400, getStatusCode(hlErr))
	assert.Equal(t, 400, getStatusCode(dErr))
	assert.Equal(t, 400, getStatusCode(sErr))
	assert.Equal(t, 400, getStatusCode(opErr))
	assert.Equal(t, 400, getStatusCode(oErr))
	assert.Equal(t, 400, getStatusCode(fErr))
}

func getStatusCode(error []error) int {
	switch e := error[0].(type) {
	case gostErrors.APIError:
		return e.GetHTTPErrorStatusCode()
	}

	return 0
}
