package entities

import "encoding/json"

// Datastream in SensorThings represents a collection of Observations from a Sensor. A physical Sensor will send its
// data (Observations) to a single Datastream. A Datastream can have only one Sensor and that Sensor can only
// sense one ObservedProperty.
type Datastream struct {
	ID                  string            `json:"@iot.id"`
	NavSelf             string            `json:"@iot.selfLink"`
	Description         string            `json:"descritption"`
	UnitOfMeasurement   map[string]string `json:"unitOfMeasurement"`
	ObservationType     string            `json:"observationType"`
	ObservedArea        map[string]string `json:"observedArea"`
	NavThings           string            `json:"Things@iot.navigationLink,omitempty"`
	NavSensors          string            `json:"Sensors@iot.navigationLink,omitempty"`
	NavObservations     string            `json:"Observations@iot.navigationLink,omitempty"`
	NavObservedProperty string            `json:"ObservedProperty@iot.navigationLink,omitempty"`
	Thing               *Thing            `json:"Thing"`
	Sensor              *Sensor           `json:"Sensor"`
	Observation         *Observation      `json:"Observation"`
	ObservedProperty    *ObservedProperty `json:"ObservedProperty"`
}

// ParseEntity tries to parse the given json byte array into the current entity
func (d *Datastream) ParseEntity(data []byte) error {
	datastream := &d
	err := json.Unmarshal(data, datastream)
	if err != nil {
		return err
	}

	return nil
}
