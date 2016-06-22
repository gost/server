package main

import (
	"fmt"
	"log"
	"testing"

	net "net/http"

	"github.com/geodan/gost/src/configuration"
	"github.com/geodan/gost/src/database/postgis"
	"github.com/geodan/gost/src/http"
	"github.com/geodan/gost/src/mqtt"
	"github.com/geodan/gost/src/sensorthings/api"
	"github.com/stretchr/testify/assert"
)

func TestVersionHandler(t *testing.T) {
	log.Println("do test")
	//arrange
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	// database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100)
	database := postgis.NewDatabase("192.168.40.10", 5432, "postgres", "postgres", "gost", "v1", false, 50, 100)
	api := api.NewAPI(database, cfg, mqttServer)

	gostServer := http.CreateServer("localhost", 8088, &api)
	gostServer.Start()
	versionUrl := fmt.Sprintf("%s/Version", "http://localhost:8088")

	// act
	request, _ := net.NewRequest("GET", versionUrl, nil)
	res, _ := net.DefaultClient.Do(request)

	//assert
	assert.Equal(t, 200, res.StatusCode, "result should be http 200")

	// teardown
	gostServer.Stop()
}
