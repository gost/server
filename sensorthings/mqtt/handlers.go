package mqtt

import (
	"fmt"
	"github.com/geodan/gost/sensorthings/models"
)

func ObservationHandler(a *models.API, clientID string, topic string, message []byte) {
	fmt.Printf("TOPIC: %s\n", topic)
	fmt.Printf("MSG: %s\n", message)
}
