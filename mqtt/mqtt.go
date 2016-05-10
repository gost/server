package mqtt

import (
	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/sensorthings/models"
)

// MQTT is the implementation of the MQTT server
type MQTT struct {
}

// CreateMQTTClient creates a new MQTT client
func CreateMQTTClient(config configuration.MQTTConfig) models.MQTTClient {
	return &MQTT{}
}

// Start running the MQTT client
func (m *MQTT) Start() {

}

// Stop the MQTT client
func (m *MQTT) Stop() {

}
