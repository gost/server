package entities

import (
	"encoding/json"
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
)

// Thing in SensorThings represents a physical object in the real world. A Thing is a good starting
// point in which to start creating the SensorThings model structure. A Thing has a Location and one or
// more Datastreams to collect Observations. A minimal Thing can be created without a Location and Datastream
// and there are options to create a Things with a nested linked Location and Datastream.
type Thing struct {
	BaseEntity
	Description            string                 `json:"description,omitempty"`
	Properties             map[string]interface{} `json:"properties,omitempty"`
	NavLocations           string                 `json:"Locations@iot.navigationLink,omitempty"`
	NavDatastreams         string                 `json:"Datastreams@iot.navigationLink,omitempty"`
	NavHistoricalLocations string                 `json:"HistoricalLocations@iot.navigationLink,omitempty"`
	Locations              []*Location            `json:"Locations,omitempty"`
	Datastreams            []*Datastream          `json:"Datastreams,omitempty"`
	HistoricalLocations    []*HistoricalLocation  `json:"HistoricalLocations,omitempty"`
}

// GetEntityType returns the EntityType for Thing
func (t Thing) GetEntityType() EntityType {
	return EntityTypeThing
}

// GetPropertyNames returns the available properties for a Thing
func (t *Thing) GetPropertyNames() []string {
	return []string{"id", "description", "properties"}
}

// ParseEntity tries to parse the given json byte array into the current entity
func (t *Thing) ParseEntity(data []byte) error {
	thing := &t
	err := json.Unmarshal(data, thing)
	if err != nil {
		return gostErrors.NewBadRequestError(errors.New("Unable to parse Thing"))
	}

	return nil
}

// ContainsMandatoryParams checks if all mandatory params for Thing are available before posting.
func (t *Thing) ContainsMandatoryParams() (bool, []error) {
	err := []error{}
	CheckMandatoryParam(&err, t.Description, t.GetEntityType(), "description")

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetAllLinks sets the self link and relational links
func (t *Thing) SetAllLinks(externalURL string) {
	t.SetSelfLink(externalURL)
	t.SetLinks(externalURL)
}

// SetSelfLink sets the self link for the entity
func (t *Thing) SetSelfLink(externalURL string) {
	t.NavSelf = CreateEntitySelfLink(externalURL, EntityLinkThings.ToString(), t.ID)
}

// SetLinks sets the entity specific navigation links, empty string if linked(expanded) data is not nil
func (t *Thing) SetLinks(externalURL string) {
	t.NavLocations = CreateEntityLink(t.Locations == nil, externalURL, EntityLinkThings.ToString(), EntityLinkLocations.ToString(), t.ID)
	t.NavDatastreams = CreateEntityLink(t.Datastreams == nil, externalURL, EntityLinkThings.ToString(), EntityLinkDatastreams.ToString(), t.ID)
	t.NavHistoricalLocations = CreateEntityLink(t.HistoricalLocations == nil, externalURL, EntityLinkThings.ToString(), EntityLinkHistoricalLocations.ToString(), t.ID)
}

// GetSupportedEncoding returns the supported encoding tye for this entity
func (t Thing) GetSupportedEncoding() map[int]EncodingType {
	return map[int]EncodingType{}
}
