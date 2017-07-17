package entities

import (
	"encoding/json"
	"errors"

	gostErrors "github.com/gost/server/errors"
)

// ObservedProperty in SensorThings represents the physical phenomenon being observed by the Sensor. An ObserveProperty is
// linked to a Datastream which can only have one ObserveProperty
type ObservedProperty struct {
	BaseEntity
	Name           string        `json:"name,omitempty"`
	Description    string        `json:"description,omitempty"`
	Definition     string        `json:"definition,omitempty"`
	NavDatastreams string        `json:"Datastreams@iot.navigationLink,omitempty"`
	Datastreams    []*Datastream `json:"Datastreams,omitempty"`
}

// GetEntityType returns the EntityType for ObservedProperty
func (o ObservedProperty) GetEntityType() EntityType {
	return EntityTypeObservedProperty
}

// GetPropertyNames returns the available properties for a ObservedProperty
func (o *ObservedProperty) GetPropertyNames() []string {
	return []string{"id", "name", "description", "definition"}
}

// ParseEntity tries to parse the given json byte array into the current entity
func (o *ObservedProperty) ParseEntity(data []byte) error {
	op := &o
	err := json.Unmarshal(data, op)
	if err != nil {
		return gostErrors.NewBadRequestError(errors.New("Unable to parse ObservedProperty"))
	}

	return nil
}

// ContainsMandatoryParams checks if all mandatory params for ObservedProperty are available before posting.
func (o *ObservedProperty) ContainsMandatoryParams() (bool, []error) {
	err := []error{}
	CheckMandatoryParam(&err, o.Name, o.GetEntityType(), "name")
	CheckMandatoryParam(&err, o.Definition, o.GetEntityType(), "definition")
	CheckMandatoryParam(&err, o.Description, o.GetEntityType(), "description")

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetAllLinks sets the self link and relational links
func (o *ObservedProperty) SetAllLinks(externalURL string) {
	o.SetSelfLink(externalURL)
	o.SetLinks(externalURL)

	for _, d := range o.Datastreams {
		d.SetAllLinks(externalURL)
	}
}

// SetSelfLink sets the self link for the entity
func (o *ObservedProperty) SetSelfLink(externalURL string) {
	o.NavSelf = CreateEntitySelfLink(externalURL, EntityLinkObservedProperties.ToString(), o.ID)
}

// SetLinks sets the entity specific navigation links, empty string if linked(expanded) data is not nil
func (o *ObservedProperty) SetLinks(externalURL string) {
	o.NavDatastreams = CreateEntityLink(o.Datastreams == nil, externalURL, EntityLinkObservedProperties.ToString(), EntityLinkDatastreams.ToString(), o.ID)
}

// GetSupportedEncoding returns the supported encoding tye for this entity
func (o ObservedProperty) GetSupportedEncoding() map[int]EncodingType {
	return map[int]EncodingType{}
}
