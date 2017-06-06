package api

import (
	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/database/postgis"
	"github.com/geodan/gost/mqtt"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/odata"

	"github.com/gost/godata"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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
