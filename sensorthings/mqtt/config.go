package mqtt

import (
	"fmt"

	"github.com/gost/server/sensorthings/models"
)

// CreateTopics creates the pre-defined MQTT Topics
func CreateTopics(prefix string) []models.Topic {
	topics := []models.Topic{
		{
			Path:    fmt.Sprintf("%s/#", prefix),
			Handler: MainMqttHandler,
		},
	}

	return topics
}
