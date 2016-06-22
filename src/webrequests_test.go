package main

import (
	"fmt"
	"testing"

	"github.com/geodan/gost/src/configuration"
	"github.com/geodan/gost/src/database/postgis"
	"github.com/geodan/gost/src/http"
	"github.com/geodan/gost/src/mqtt"
	"github.com/geodan/gost/src/sensorthings/api"
	"github.com/stretchr/testify/assert"
)

func TestVersionHandler(t *testing.T) {
	//arrange
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100)
	a := api.NewAPI(database, cfg, mqttServer)
	gostServer := http.CreateServer(a.GetConfig().Server.Host, a.GetConfig().Server.Port, api)
	gostServer.Start()
	versionUrl := fmt.Sprintf("%s/Version", gostserver.URL)

	// act
	request, err := net.http.NewRequest("GET", versionUrl, nil)
	res, err := http.DefaultClient.Do(request)

	//assert
	assert.Equal(t, 200, res.StatusCode, "result should be http 200")

	// teardown
	gostServer.Stop()
}
