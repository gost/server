package entities

import "encoding/json"

// Observation in SensorThings represents a single Sensor reading of an ObservedProperty. A physical device, a Sensor, sends
// Observations to a specified Datastream. An Observation requires a FeaturOfInterest entity, if none is provided in the request,
// the Location of the Thing associated with the Datastream, will be assigned to the new Observation as the FeaturOfInterest.
type Observation struct {
	ID                   string             `json:"@iot.id"`
	NavSelf              string             `json:"@iot.selfLink"`
	PhenomenonTime       string             `json:"phenomenonTime"`
	Result               string             `json:"result"`
	ResultTime           string             `json:"resultTime"`
	ResultQuality        string             `json:"resultQuality"`
	ValidTime            string             `json:"validTime"`
	Parameters           string             `json:"parameters"`
	NavDatastreams       string             `json:"Datastreams@iot.navigationLink,omitempty"`
	NavFeatureOfInterest string             `json:"FeatureOfInterest@iot.navigationLink,omitempty"`
	Datastreams          []*Datastream      `json:"Datastreams,omitempty"`
	FeatureOfInterest    *FeatureOfInterest `json:"FeatureOfInterest,omitempty"`
}

// ParseEntity tries to parse the given json byte array into the current entity
func (o *Observation) ParseEntity(data []byte) error {
	observation := &o
	err := json.Unmarshal(data, observation)
	if err != nil {
		return err
	}

	return nil
}
