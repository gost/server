package entities

import (
	"encoding/json"

	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"time"
)

// Observation in SensorThings represents a single Sensor reading of an ObservedProperty. A physical device, a Sensor, sends
// Observations to a specified Datastream. An Observation requires a FeaturOfInterest entity, if none is provided in the request,
// the Location of the Thing associated with the Datastream, will be assigned to the new Observation as the FeaturOfInterest.
type Observation struct {
	BaseEntity
	NavSelf              string                 `json:"@iot.selfLink,omitempty"`
	PhenomenonTime       string                 `json:"phenomenonTime,omitempty"`
	Result               interface{}            `json:"result,omitempty"`
	ResultTime           string                 `json:"resultTime,omitempty"`
	ResultQuality        string                 `json:"resultQuality,omitempty"`
	ValidTime            string                 `json:"validTime,omitempty"`
	Parameters           map[string]interface{} `json:"parameters,omitempty"`
	NavDatastream        string                 `json:"Datastream@iot.navigationLink,omitempty"`
	NavFeatureOfInterest string                 `json:"FeatureOfInterest@iot.navigationLink,omitempty"`
	Datastream           *Datastream            `json:"Datastream,omitempty"`
	FeatureOfInterest    *FeatureOfInterest     `json:"FeatureOfInterest,omitempty"`
}

// GetEntityType returns the EntityType for Observation
func (o Observation) GetEntityType() EntityType {
	return EntityTypeObservation
}

// ParseEntity tries to parse the given json byte array into the current entity
func (o *Observation) ParseEntity(data []byte) error {
	observation := &o
	err := json.Unmarshal(data, observation)
	if err != nil {
		return gostErrors.NewBadRequestError(errors.New("Unable to parse Observation"))
	}

	return nil
}

// ContainsMandatoryParams checks if all mandatory params for Observation are available before posting.
func (o *Observation) ContainsMandatoryParams() (bool, []error) {
	// When a SensorThings service receives a POST Observations without phenonmenonTime, the service SHALL
	// assign the current server time to the value of the phenomenonTime.
	if len(o.PhenomenonTime) == 0 {
		o.PhenomenonTime = time.Now().UTC().Format(time.RFC3339Nano)
	}

	// When a SensorThings service receives a POST Observations without resultTime, the service SHALL assign a
	// null value to the resultTime.
	if len(o.ResultTime) == 0 {
		o.ResultTime = "NULL"
	}

	var err []error
	CheckMandatoryParam(&err, o.PhenomenonTime, o.GetEntityType(), "phenomenonTime")
	CheckMandatoryParam(&err, o.Result, o.GetEntityType(), "result")
	CheckMandatoryParam(&err, o.ResultTime, o.GetEntityType(), "resultTime")
	CheckMandatoryParam(&err, o.Datastream, o.GetEntityType(), "Datastream")

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetLinks sets the entity specific navigation links if needed
func (o *Observation) SetLinks(externalURL string) {
	o.NavSelf = CreateEntitySelfLink(externalURL, EntityLinkObservations.ToString(), o.ID)
	o.NavDatastream = CreateEntityLink(o.Datastream == nil, externalURL, EntityLinkObservations.ToString(), EntityTypeDatastream.ToString(), o.ID)
	o.NavFeatureOfInterest = CreateEntityLink(o.FeatureOfInterest == nil, externalURL, EntityLinkObservations.ToString(), EntityTypeFeatureOfInterest.ToString(), o.ID)
}

// MarshalPostgresJSON marshalls an observation entity for saving into PostgreSQL
func (o Observation) MarshalPostgresJSON() ([]byte, error) {
	return json.Marshal(&struct {
		PhenomenonTime string                 `json:"phenomenonTime,omitempty"`
		Result         interface{}            `json:"result,omitempty"`
		ResultTime     string                 `json:"resultTime,omitempty"`
		ResultQuality  string                 `json:"resultQuality,omitempty"`
		ValidTime      string                 `json:"validTime,omitempty"`
		Parameters     map[string]interface{} `json:"parameters,omitempty"`
	}{
		PhenomenonTime: o.PhenomenonTime,
		Result:         o.Result,
		ResultTime:     o.ResultTime,
		ResultQuality:  o.ResultQuality,
		ValidTime:      o.ValidTime,
		Parameters:     o.Parameters,
	})
}

// GetSupportedEncoding returns the supported encoding tye for this entity
func (o Observation) GetSupportedEncoding() map[int]EncodingType {
	return map[int]EncodingType{}
}
