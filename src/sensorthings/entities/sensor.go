package entities

import (
	"encoding/json"
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
)

// Sensor in SensorThings represents the physical device capable of observing a physical property and converting
// it to an electrical impulse and be converted to a empirical value to represent a measurement value of the physical property
type Sensor struct {
	ID             string        `json:"@iot.id,omitempty"`
	NavSelf        string        `json:"@iot.selfLink,omitempty"`
	Description    string        `json:"description,omitempty"`
	EncodingType   string        `json:"encodingtype,omitempty"`
	Metadata       string        `json:"metadata,omitempty"`
	NavDatastreams string        `json:"Datastreams@iot.navigationLink,omitempty"`
	Datastreams    []*Datastream `json:"Datastreams,omitempty"`
}

// GetEntityType returns the EntityType for Sensor
func (s Sensor) GetEntityType() EntityType {
	return EntityTypeSensor
}

// ParseEntity tries to parse the given json byte array into the current entity
func (s *Sensor) ParseEntity(data []byte) error {
	sensor := &s
	err := json.Unmarshal(data, sensor)
	if err != nil {
		return gostErrors.NewBadRequestError(errors.New("Unable to parse Sensor"))
	}

	return nil
}

// ContainsMandatoryParams checks if all mandatory params for Sensor are available before posting.
func (s Sensor) ContainsMandatoryParams() (bool, []error) {
	err := []error{}
	CheckMandatoryParam(&err, s.Description, s.GetEntityType(), "description")
	CheckMandatoryParam(&err, s.EncodingType, s.GetEntityType(), "encodingtype")
	CheckMandatoryParam(&err, s.Metadata, s.GetEntityType(), "metadata")

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetLinks sets the entity specific navigation links if needed
func (s *Sensor) SetLinks(externalURL string) {
	s.NavSelf = CreateEntitySelfLink(externalURL, EntityLinkSensors.ToString(), s.ID)
	s.NavDatastreams = CreateEntityLink(s.Datastreams == nil, externalURL, EntityLinkSensors.ToString(), EntityLinkDatastreams.ToString(), s.ID)
}

// GetSupportedEncoding returns the supported encoding tye for this entity
func (s Sensor) GetSupportedEncoding() map[int]EncodingType {
	return map[int]EncodingType{EncodingSensorML.Code: EncodingSensorML, EncodingPDF.Code: EncodingPDF}
}
