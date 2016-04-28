package entities

import "encoding/json"

// ObservedProperty in SensorThings represents the physical phenomenon being observed by the Sensor. An ObserveProperty is
// linked to a Datatream which can only have one ObserveProperty
type ObservedProperty struct {
	ID             string        `json:"@iot.id"`
	NavSelf        string        `json:"@iot.selfLink"`
	Description    string        `json:"descritption"`
	Name           string        `json:"name"`
	Definition     string        `json:"definition"`
	NavDatastreams string        `json:"Datastreams@iot.navigationLink,omitempty"`
	Datastreams    []*Datastream `json:"Datastreams,omitempty"`
}

// ParseEntity tries to parse the given json byte array into the current entity
func (o *ObservedProperty) ParseEntity(data []byte) error {
	op := &o
	err := json.Unmarshal(data, op)
	if err != nil {
		return err
	}

	return nil
}
