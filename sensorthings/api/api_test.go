package api

import (
    "github.com/geodan/gost/configuration"
    "github.com/geodan/gost/mqtt"
    "github.com/geodan/gost/database/postgis"
    
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateApi(t *testing.T) {
    // arrange
	cfg := configuration.Config{}
    mqttServer := mqtt.NewMQTTServer(configuration.MQTTConfig{})
    database := postgis.NewDatabase("",123,"", "", "", "",false)
    stAPI := NewAPI(database, cfg, mqttServer)
    
    // act
    ep := stAPI.GetEndpoints()
    
    // assert
	assert.NotNil(t,ep)
}