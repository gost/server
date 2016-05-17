package api

import (
	"github.com/geodan/gost/src/configuration"
	"github.com/geodan/gost/src/database/postgis"
	"github.com/geodan/gost/src/mqtt"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateApi(t *testing.T) {
	// arrange
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100)
	stAPI := NewAPI(database, cfg, mqttServer)

	// act
	ep := stAPI.GetEndpoints()
	endpoints := *ep

	// assert
	assert.NotNil(t, ep)
	assert.NotEqual(t, len(endpoints), 0, "Endpoints empty")
}
