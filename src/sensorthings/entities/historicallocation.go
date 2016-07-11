package entities

import (
	"encoding/json"
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"time"
)

// HistoricalLocation in sensorthings represents the current and previous locations of a thing including time
type HistoricalLocation struct {
	BaseEntity
	Time         string      `json:"time,omitempty"`
	NavThing     string      `json:"Thing@iot.navigationLink,omitempty"`
	NavLocations string      `json:"Locations@iot.navigationLink,omitempty"`
	Thing        *Thing      `json:"Thing,omitempty"`
	Locations    []*Location `json:"Locations,omitempty"`
}

// GetEntityType returns the EntityType for HistoricalLocation
func (h HistoricalLocation) GetEntityType() EntityType {
	return EntityTypeHistoricalLocation
}

// GetPropertyNames returns the available properties for a HistoricalLocation
func (h *HistoricalLocation) GetPropertyNames() []string {
	return []string{"id", "time"}
}

// ParseEntity tries to parse the given json byte array into the current entity
func (h *HistoricalLocation) ParseEntity(data []byte) error {
	hl := &h
	err := json.Unmarshal(data, hl)
	if err != nil {
		return gostErrors.NewBadRequestError(errors.New("Unable to parse HistoricalLocation"))
	}

	return nil
}

// ContainsMandatoryParams checks if all mandatory params for a HistoricalLocation are available before posting
func (h *HistoricalLocation) ContainsMandatoryParams() (bool, []error) {
	err := []error{}

	if len(h.Time) == 0 {
		h.Time = time.Now().UTC().Format(time.RFC3339Nano)
	}

	CheckMandatoryParam(&err, h.Time, h.GetEntityType(), "time")
	CheckMandatoryParam(&err, h.Thing, h.GetEntityType(), "Thing")
	if h.Thing != nil {
		CheckMandatoryParam(&err, h.Thing.ID, h.GetEntityType(), "Thing.ID")
	}

	CheckMandatoryParam(&err, h.Locations, h.GetEntityType(), "Location")
	if len(h.Locations) == 0 {
		err = append(err, gostErrors.NewBadRequestError(errors.New("Missing location")))
	} else {
		CheckMandatoryParam(&err, h.Locations[0].ID, h.GetEntityType(), "Location.ID")
	}

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetAllLinks sets the self link and relational links
func (h *HistoricalLocation) SetAllLinks(externalURL string) {
	h.SetSelfLink(externalURL)
	h.SetLinks(externalURL)
}

// SetSelfLink sets the self link for the entity
func (h *HistoricalLocation) SetSelfLink(externalURL string) {
	h.NavSelf = CreateEntitySelfLink(externalURL, EntityLinkHistoricalLocations.ToString(), h.ID)
}

// SetLinks sets the entity specific navigation links, empty string if linked(expanded) data is not nil
func (h *HistoricalLocation) SetLinks(externalURL string) {
	h.NavThing = CreateEntityLink(h.Thing == nil, externalURL, EntityLinkHistoricalLocations.ToString(), EntityTypeThing.ToString(), h.ID)
	h.NavLocations = CreateEntityLink(h.Locations == nil, externalURL, EntityLinkHistoricalLocations.ToString(), EntityLinkLocations.ToString(), h.ID)
}

// GetSupportedEncoding returns the supported encoding tye for this entity
func (h HistoricalLocation) GetSupportedEncoding() map[int]EncodingType {
	return map[int]EncodingType{}
}
