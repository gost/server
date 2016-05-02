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
	NavDatastream        string             `json:"Datastream@iot.navigationLink,omitempty"`
	NavFeatureOfInterest string             `json:"FeatureOfInterest@iot.navigationLink,omitempty"`
	Datastream           *Datastream        `json:"Datastream,omitempty"`
	FeatureOfInterest    *FeatureOfInterest `json:"FeatureOfInterest,omitempty"`
}

// GetEntityType returns the EntityType for Observation
func (o *Observation) GetEntityType() EntityType {
	return EntityTypeObservation
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

// ContainsMandatoryParams checks if all mandatory params for Observation are available before posting.
func (o *Observation) ContainsMandatoryParams() (bool, []error) {
	err := []error{}
	CheckMandatoryParam(&err, o.PhenomenonTime, o.GetEntityType(), "phenomenonTime")
	CheckMandatoryParam(&err, o.Result, o.GetEntityType(), "result")
	CheckMandatoryParam(&err, o.ResultTime, o.GetEntityType(), "resultTime")

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetLinks sets the entity specific navigation links if needed
func (o *Observation) SetLinks(externalURL string) {
	o.NavSelf = CreateEntitySefLink(externalURL, EntityLinkObservations.ToString(), o.ID)
	o.NavDatastream = CreateEntityLink(o.Datastream == nil, EntityLinkObservations.ToString(), EntityTypeDatastream.ToString(), o.ID)
	o.NavFeatureOfInterest = CreateEntityLink(o.FeatureOfInterest == nil, EntityLinkObservations.ToString(), EntityTypeFeatureOfInterest.ToString(), o.ID)
}
