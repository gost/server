package entities

import "encoding/json"

// FeatureOfInterest in SensorThings represents the phenomena an Observation is detecting. In some cases a FeatureOfInterest
// can be the Location of the Sensor and therefore of the Observation. A FeatureOfInterest is linked to a single Observation
type FeatureOfInterest struct {
	ID              string       `json:"@iot.id"`
	NavSelf         string       `json:"@iot.selfLink"`
	Description     string       `json:"descritption"`
	EncodingType    string       `json:"encodingtype"`
	Feature         string       `json:"feature"`
	NavObservations string       `json:"Observations@iot.navigationLink,omitempty"`
	Observation     *Observation `json:"Observation"`
}

// ParseEntity tries to parse the given json byte array into the current entity
func (f *FeatureOfInterest) ParseEntity(data []byte) error {
	foi := &f
	err := json.Unmarshal(data, foi)
	if err != nil {
		return err
	}

	return nil
}
