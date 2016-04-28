package entities

import "encoding/json"

// Sensor in SensorThings represents the physical device capable of observing a physical property and converting
// it to an electrical impulse and be converted to a empirical value to represent a measurement value of the physical property
type Sensor struct {
	ID             string        `json:"@iot.id"`
	NavSelf        string        `json:"@iot.selfLink"`
	Description    string        `json:"descritption"`
	EncodingType   string        `json:"encodingtype"`
	Metadata       string        `json:"metadata"`
	NavDatastreams string        `json:"Datastreams@iot.navigationLink,omitempty"`
	Datastreams    []*Datastream `json:"Datastreams,omitempty"`
}

// ParseEntity tries to parse the given json byte array into the current entity
func (s *Sensor) ParseEntity(data []byte) error {
	sensor := &s
	err := json.Unmarshal(data, sensor)
	if err != nil {
		return err
	}

	return nil
}
