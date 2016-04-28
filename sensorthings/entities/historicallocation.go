package entities

import "encoding/json"

// HistoricalLocation in sensorthings represents the current and previous locations of a thing including time
type HistoricalLocation struct {
	ID           string      `json:"@iot.id"`
	NavSelf      string      `json:"@iot.selfLink"`
	Time         string      `json:"@iot.selfLink"`
	NavThing     string      `json:"Thing@iot.navigationLink,omitempty"`
	NavLocations string      `json:"Locations@iot.navigationLink,omitempty"`
	Thing        *Thing      `json:"Thing"`
	Locations    []*Location `json:"Locations,omitempty"`
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
