package entities

import (
	"encoding/json"
	"github.com/geodan/gost/errors"
)

// Datastream in SensorThings represents a collection of Observations from a Sensor. A physical Sensor will send its
// data (Observations) to a single Datastream. A Datastream can have only one Sensor and that Sensor can only
// sense one ObservedProperty.
type Datastream struct {
	ID                  string            `json:"@iot.id"`
	NavSelf             string            `json:"@iot.selfLink"`
	Description         string            `json:"description"`
	UnitOfMeasurement   map[string]string `json:"unitOfMeasurement"`
	ObservationType     string            `json:"observationType"`
	ObservedArea        map[string]string `json:"observedArea"`
	NavThing            string            `json:"Thing@iot.navigationLink,omitempty"`
	NavSensor           string            `json:"Sensor@iot.navigationLink,omitempty"`
	NavObservations     string            `json:"Observations@iot.navigationLink,omitempty"`
	NavObservedProperty string            `json:"ObservedProperty@iot.navigationLink,omitempty"`
	Thing               *Thing            `json:"Thing,omitempty"`
	Sensor              *Sensor           `json:"Sensor,omitempty"`
	Observations        *[]Observation    `json:"Observations,omitempty"`
	ObservedProperty    *ObservedProperty `json:"ObservedProperty,omitempty"`
}

// GetEntityType returns the EntityType for Datastream
func (d *Datastream) GetEntityType() EntityType {
	return EntityTypeDatastream
}

// ParseEntity tries to parse the given json byte array into the current entity
func (d *Datastream) ParseEntity(data []byte) error {
	datastream := &d
	err := json.Unmarshal(data, datastream)
	if err != nil {
		return errors.NewBadRequestError(err)
	}

	return nil
}

// ContainsMandatoryParams checks if all mandatory params for a Datastream are available before posting
func (d *Datastream) ContainsMandatoryParams() (bool, []error) {
	err := []error{}
	CheckMandatoryParam(&err, d.Description, d.GetEntityType(), "description")
	CheckMandatoryParam(&err, d.UnitOfMeasurement, d.GetEntityType(), "unitOfMeasurement")
	CheckMandatoryParam(&err, d.ObservationType, d.GetEntityType(), "observationType")

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetLinks sets the entity specific navigation links, empty string if linked(expanded) data is not nil
func (d *Datastream) SetLinks(externalURL string) {
	d.NavSelf = CreateEntitySefLink(externalURL, EntityLinkDatastreams.ToString(), d.ID)
	d.NavThing = CreateEntityLink(d.Thing == nil, EntityLinkDatastreams.ToString(), EntityTypeThing.ToString(), d.ID)
	d.NavSensor = CreateEntityLink(d.Sensor == nil, EntityLinkDatastreams.ToString(), EntityTypeSensor.ToString(), d.ID)
	d.NavObservations = CreateEntityLink(d.Observations == nil, EntityLinkDatastreams.ToString(), EntityLinkObservations.ToString(), d.ID)
	d.NavObservedProperty = CreateEntityLink(d.ObservedProperty == nil, EntityLinkDatastreams.ToString(), EntityTypeObservedProperty.ToString(), d.ID)
}
