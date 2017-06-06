package entities

import (
	"encoding/json"
	"errors"

	gostErrors "github.com/geodan/gost/errors"
)

// Location entity locates the Thing or the Things it associated with. A Thing’s Location entity is
// defined as the last known location of the Thing.
// A Thing’s Location may be identical to the Thing’s Observations’ FeatureOfInterest. In the context of the IoT,
// the principle location of interest is usually associated with the location of the Thing, especially for in-situ
// sensing applications. For example, the location of interest of a wifi-connected thermostat should be the building
// or the room in which the smart thermostat is located. And the FeatureOfInterest of the Observations made by the
// thermostat (e.g., room temperature readings) should also be the building or the room. In this case, the content
// of the smart thermostat’s location should be the same as the content of the temperature readings’ feature of interest.
type Location struct {
	BaseEntity
	Name                   string                 `json:"name,omitempty"`
	Description            string                 `json:"description,omitempty"`
	EncodingType           string                 `json:"encodingType,omitempty"`
	Location               map[string]interface{} `json:"location,omitempty"`
	NavThings              string                 `json:"Things@iot.navigationLink,omitempty"`
	NavHistoricalLocations string                 `json:"HistoricalLocations@iot.navigationLink,omitempty"`
	Things                 []*Thing               `json:"Things,omitempty"`
	HistoricalLocations    []*HistoricalLocation  `json:"HistoricalLocations,omitempty"`
}

// GetEntityType returns the EntityType for Location
func (l Location) GetEntityType() EntityType {
	return EntityTypeLocation
}

// GetPropertyNames returns the available properties for a Location
func (l *Location) GetPropertyNames() []string {
	return []string{"id", "name", "description", "encodingType", "location"}
}

// ParseEntity tries to parse the given json byte array into the current entity
func (l *Location) ParseEntity(data []byte) error {
	location := &l
	err := json.Unmarshal(data, location)
	if err != nil {
		return gostErrors.NewBadRequestError(errors.New("Unable to parse Location"))
	}

	return nil
}

// ContainsMandatoryParams checks if all mandatory params for Location are available before posting.
func (l *Location) ContainsMandatoryParams() (bool, []error) {
	err := []error{}
	CheckMandatoryParam(&err, l.Name, l.GetEntityType(), "name")
	CheckMandatoryParam(&err, l.Description, l.GetEntityType(), "description")
	CheckMandatoryParam(&err, l.EncodingType, l.GetEntityType(), "encodingType")
	CheckMandatoryParam(&err, l.Location, l.GetEntityType(), "location")

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetAllLinks sets the self link and relational links
func (l *Location) SetAllLinks(externalURL string) {
	l.SetSelfLink(externalURL)
	l.SetLinks(externalURL)

	for _, t := range l.Things {
		t.SetAllLinks(externalURL)
	}

	for _, hl := range l.HistoricalLocations {
		hl.SetAllLinks(externalURL)
	}
}

// SetSelfLink sets the self link for the entity
func (l *Location) SetSelfLink(externalURL string) {
	l.NavSelf = CreateEntitySelfLink(externalURL, EntityLinkLocations.ToString(), l.ID)
}

// SetLinks sets the entity specific navigation links, empty string if linked(expanded) data is not nil
func (l *Location) SetLinks(externalURL string) {
	l.NavThings = CreateEntityLink(l.Things == nil, externalURL, EntityLinkLocations.ToString(), EntityLinkThings.ToString(), l.ID)
	l.NavHistoricalLocations = CreateEntityLink(l.HistoricalLocations == nil, externalURL, EntityLinkLocations.ToString(), EntityLinkHistoricalLocations.ToString(), l.ID)
}
