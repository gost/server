package mqtt

import (
	"github.com/geodan/gost/sensorthings/models"
)

// CreateTopics creates the pre-defined MQTT Topics
func CreateTopics() []models.Topic {
	topics := []models.Topic{
		models.Topic{
			Path:    "Datastreams(1)/Observations",
			Handler: ObservationHandler,
		},
	}

	return topics
}
