package mqtt

import (
	"github.com/gost/server/sensorthings/models"
)

// CreateTopics creates the pre-defined MQTT Topics
func CreateTopics() []models.Topic {
	topics := []models.Topic{
		{
			Path:    "GOST/#",
			Handler: MainMqttHandler,
		},
	}

	return topics
}
