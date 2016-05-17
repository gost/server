package mqtt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/geodan/gost/src/configuration"
)

func TestMqtt(t *testing.T) {
	// arrange
	config := configuration.MQTTConfig{}
	config.Host = "iot.eclipse.org"
	config.Port = 1883

	// act
	mqttClient := CreateMQTTClient(config)
	assert.NotNil(t, mqttClient, "function should return MqqtClient")
}
