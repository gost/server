package entities

import (
	"encoding/json"
	"errors"

	gostErrors "github.com/geodan/gost/src/errors"
)

// Sensor in SensorThings represents the physical device capable of observing a physical property and converting
// it to an electrical impulse and be converted to a empirical value to represent a measurement value of the physical property
type Sensor struct {
	BaseEntity
	Name           string        `json:"name,omitempty"`
	Description    string        `json:"description,omitempty"`
	EncodingType   string        `json:"encodingType,omitempty"`
	Metadata       string        `json:"metadata,omitempty"`
	NavDatastreams string        `json:"Datastreams@iot.navigationLink,omitempty"`
	Datastreams    []*Datastream `json:"Datastreams,omitempty"`
}

// GetEntityType returns the EntityType for Sensor
func (s Sensor) GetEntityType() EntityType {
	return EntityTypeSensor
}

// GetPropertyNames returns the available properties for a Sensor
func (s *Sensor) GetPropertyNames() []string {
	return []string{"id", "name", "description", "encodingType", "metadata"}
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
func (s *Sensor) ContainsMandatoryParams() (bool, []error) {
	err := []error{}
	CheckMandatoryParam(&err, s.Name, s.GetEntityType(), "name")
	CheckMandatoryParam(&err, s.Description, s.GetEntityType(), "description")
	CheckMandatoryParam(&err, s.EncodingType, s.GetEntityType(), "encodingType")
	CheckMandatoryParam(&err, s.Metadata, s.GetEntityType(), "metadata")

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetAllLinks sets the self link and relational links
func (s *Sensor) SetAllLinks(externalURL string) {
	s.SetSelfLink(externalURL)
	s.SetLinks(externalURL)
}

// SetSelfLink sets the self link for the entity
func (s *Sensor) SetSelfLink(externalURL string) {
	s.NavSelf = CreateEntitySelfLink(externalURL, EntityLinkSensors.ToString(), s.ID)
}

// SetLinks sets the entity specific navigation links, empty string if linked(expanded) data is not nil
func (s *Sensor) SetLinks(externalURL string) {
	s.NavDatastreams = CreateEntityLink(s.Datastreams == nil, externalURL, EntityLinkSensors.ToString(), EntityLinkDatastreams.ToString(), s.ID)
}

// GetSupportedEncoding returns the supported encoding tye for this entity
func (s Sensor) GetSupportedEncoding() map[int]EncodingType {
	return map[int]EncodingType{EncodingSensorML.Code: EncodingSensorML, EncodingPDF.Code: EncodingPDF, EncodingTextHTML.Code: EncodingTextHTML, EncodingTypeDescription.Code: EncodingTypeDescription}
}
