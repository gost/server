package entities

import "encoding/json"

// HistoricalLocation in sensorthings represents the current and previous locations of a thing including time
type HistoricalLocation struct {
	ID           string      `json:"@iot.id,omitempty"`
	NavSelf      string      `json:"@iot.selfLink,omitempty"`
	Time         string      `json:"time,omitempty"`
	NavThing     string      `json:"Thing@iot.navigationLink,omitempty"`
	NavLocations string      `json:"Locations@iot.navigationLink,omitempty"`
	Thing        *Thing      `json:"Thing,omitempty"`
	Locations    []*Location `json:"Locations,omitempty"`
}

// GetEntityType returns the EntityType for HistoricalLocation
func (h *HistoricalLocation) GetEntityType() EntityType {
	return EntityTypeHistoricalLocation
}

// ParseEntity tries to parse the given json byte array into the current entity
func (h *HistoricalLocation) ParseEntity(data []byte) error {
	hl := &h
	err := json.Unmarshal(data, hl)
	if err != nil {
		return err
	}

	return nil
}

// ContainsMandatoryParams checks if all mandatory params for a HistoricalLocation are available before posting
func (h *HistoricalLocation) ContainsMandatoryParams() (bool, []error) {
	err := []error{}
	CheckMandatoryParam(&err, h.Time, h.GetEntityType(), "time")
	CheckMandatoryParam(&err, h.Thing, h.GetEntityType(), "Thing")

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetLinks sets the entity specific navigation links if needed
func (h *HistoricalLocation) SetLinks(externalURL string) {
	h.NavSelf = CreateEntitySefLink(externalURL, EntityLinkHistoricalLocations.ToString(), h.ID)
	h.NavThing = CreateEntityLink(h.Thing == nil, EntityLinkHistoricalLocations.ToString(), EntityTypeThing.ToString(), h.ID)
	h.NavLocations = CreateEntityLink(h.Locations == nil, EntityLinkHistoricalLocations.ToString(), EntityLinkLocations.ToString(), h.ID)
}
