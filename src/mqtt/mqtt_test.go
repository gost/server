package mqtt

import (
	"github.com/geodan/gost/src/configuration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMqtt(t *testing.T) {
	// arrange
	config := configuration.MQTTConfig{}
	config.Host = "iot.eclipse.org"
	config.Port = 1883

	// act
	mqttClient := CreateMQTTClient(config)
	
	// assert
	assert.NotNil(t, mqttClient, "function should return MqqtClient")
}
