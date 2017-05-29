package mqtt

import (
	"strings"

	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
)

var topics = map[string]models.MQTTInternalHandler{
	"GOST/Datastreams()/Observations": observationsByDatastream,
}

// MainMqttHandler handles all messages on GOST/# and maps them to the appropriate
// handler. Mapping is needed because of the ODATA (id) format
func MainMqttHandler(a *models.API, topic string, message []byte) {
	topicMapName := topic
	id := ""
	if strings.Contains(topic, "(") {
		i := strings.Index(topic, "(")
		i2 := strings.Index(topic, ")")
		first := topic[0 : i+1]
		id = topic[i+1 : i2]
		last := topic[i2:]
		topicMapName = first + last
	}

	h := topics[topicMapName]
	if h != nil {
		h(a, message, id)
	}
}

func observationsByDatastream(a *models.API, message []byte, id string) {
	o := entities.Observation{}
	err := o.ParseEntity(message)
	if err != nil {
		return
	}

	api := *a
	api.PostObservationByDatastream(id, &o)
}
