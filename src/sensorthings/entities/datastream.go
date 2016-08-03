package entities

import (
	"encoding/json"
	"errors"

	gostErrors "github.com/geodan/gost/src/errors"
)

// Datastream in SensorThings represents a collection of Observations from a Sensor. A physical Sensor will send its
// data (Observations) to a single Datastream. A Datastream can have only one Sensor and that Sensor can only
// sense one ObservedProperty.
type Datastream struct {
	BaseEntity
	Name                string                 `json:"name,omitempty"`
	Description         string                 `json:"description,omitempty"`
	UnitOfMeasurement   map[string]interface{} `json:"unitOfMeasurement,omitempty"`
	ObservationType     string                 `json:"observationType,omitempty"`
	ObservedArea        map[string]interface{} `json:"observedArea,omitempty"`
	NavThing            string                 `json:"Thing@iot.navigationLink,omitempty"`
	NavSensor           string                 `json:"Sensor@iot.navigationLink,omitempty"`
	NavObservations     string                 `json:"Observations@iot.navigationLink,omitempty"`
	NavObservedProperty string                 `json:"ObservedProperty@iot.navigationLink,omitempty"`
	Thing               *Thing                 `json:"Thing,omitempty"`
	Sensor              *Sensor                `json:"Sensor,omitempty"`
	Observations        []*Observation         `json:"Observations,omitempty"`
	ObservedProperty    *ObservedProperty      `json:"ObservedProperty,omitempty"`
	PhenomenonTime      string                 `json:"phenomenonTime,omitempty"`
	ResultTime          string                 `json:"resultTime,omitempty"`
}

// "phenomenonTime"

// GetEntityType returns the EntityType for Datastream
func (d Datastream) GetEntityType() EntityType {
	return EntityTypeDatastream
}

// GetPropertyNames returns the available properties for a Datastream
func (d *Datastream) GetPropertyNames() []string {
	return []string{"id", "name", "description", "unitOfMeasurement", "observationType", "observedArea", "phenomenonTime", "resultTime"}
}

// ParseEntity tries to parse the given json byte array into the current entity
func (d *Datastream) ParseEntity(data []byte) error {
	datastream := &d
	err := json.Unmarshal(data, datastream)
	if err != nil {
		return gostErrors.NewBadRequestError(errors.New("Unable to parse Datastream"))
	}

	return nil
}

// ContainsMandatoryParams checks if all mandatory params for a Datastream are available before posting
func (d *Datastream) ContainsMandatoryParams() (bool, []error) {
	err := []error{}
	CheckMandatoryParam(&err, d.Name, d.GetEntityType(), "name")
	CheckMandatoryParam(&err, d.Description, d.GetEntityType(), "description")
	CheckMandatoryParam(&err, d.UnitOfMeasurement, d.GetEntityType(), "unitOfMeasurement")
	CheckMandatoryParam(&err, d.ObservationType, d.GetEntityType(), "observationType")
	CheckMandatoryParam(&err, d.Thing, d.GetEntityType(), "Thing")
	CheckMandatoryParam(&err, d.Sensor, d.GetEntityType(), "Sensor")
	CheckMandatoryParam(&err, d.ObservedProperty, d.GetEntityType(), "ObservedProperty")

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetAllLinks sets the self link and relational links
func (d *Datastream) SetAllLinks(externalURL string) {
	d.SetSelfLink(externalURL)
	d.SetLinks(externalURL)
}

// SetSelfLink sets the self link for the entity
func (d *Datastream) SetSelfLink(externalURL string) {
	d.NavSelf = CreateEntitySelfLink(externalURL, EntityLinkDatastreams.ToString(), d.ID)
}

// SetLinks sets the entity specific navigation links, empty string if linked(expanded) data is not nil
func (d *Datastream) SetLinks(externalURL string) {
	d.NavThing = CreateEntityLink(d.Thing == nil, externalURL, EntityLinkDatastreams.ToString(), EntityTypeThing.ToString(), d.ID)
	d.NavSensor = CreateEntityLink(d.Sensor == nil, externalURL, EntityLinkDatastreams.ToString(), EntityTypeSensor.ToString(), d.ID)
	d.NavObservations = CreateEntityLink(d.Observations == nil, externalURL, EntityLinkDatastreams.ToString(), EntityLinkObservations.ToString(), d.ID)
	d.NavObservedProperty = CreateEntityLink(d.ObservedProperty == nil, externalURL, EntityLinkDatastreams.ToString(), EntityTypeObservedProperty.ToString(), d.ID)
}

// GetSupportedEncoding returns the supported encoding tye for this entity
func (d Datastream) GetSupportedEncoding() map[int]EncodingType {
	return map[int]EncodingType{}
}
