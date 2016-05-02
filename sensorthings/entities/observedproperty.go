package entities

import "encoding/json"

// ObservedProperty in SensorThings represents the physical phenomenon being observed by the Sensor. An ObserveProperty is
// linked to a Datatream which can only have one ObserveProperty
type ObservedProperty struct {
	ID             string        `json:"@iot.id"`
	NavSelf        string        `json:"@iot.selfLink"`
	Description    string        `json:"description"`
	Name           string        `json:"name"`
	Definition     string        `json:"definition"`
	NavDatastreams string        `json:"Datastreams@iot.navigationLink,omitempty"`
	Datastreams    []*Datastream `json:"Datastreams,omitempty"`
}

// GetEntityType returns the EntityType for ObservedProperty
func (o *ObservedProperty) GetEntityType() EntityType {
	return EntityTypeObservedProperty
}

// ParseEntity tries to parse the given json byte array into the current entity
func (o *ObservedProperty) ParseEntity(data []byte) error {
	op := &o
	err := json.Unmarshal(data, op)
	if err != nil {
		return err
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

// SetLinks sets the entity specific navigation links if needed
func (o *ObservedProperty) SetLinks(externalURL string) {
	o.NavSelf = CreateEntitySefLink(externalURL, EntityLinkObservedPropertys.ToString(), o.ID)
	o.NavDatastreams = CreateEntityLink(o.Datastreams == nil, EntityLinkObservedPropertys.ToString(), EntityLinkDatastreams.ToString(), o.ID)
}
