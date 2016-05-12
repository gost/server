package mqtt

import (
	"github.com/geodan/gost/src/sensorthings/models"
)

// CreateTopics creates the pre-defined MQTT Topics
func CreateTopics() []models.Topic {
	topics := []models.Topic{
		/*models.Topic{
			Path:    "$SYS/#",
			Handler: ObservationHandler,
		},*/
		models.Topic{
			Path:    "GOST/#",
			Handler: MainMqttHandler,
		},
	}

	return topics
}
