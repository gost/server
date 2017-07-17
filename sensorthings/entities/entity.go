package entities

import (
	"fmt"

	gostErrors "github.com/gost/server/errors"
	"strings"
)

// EntityType holds the name and type of a SensorThings entity.
type EntityType string

// List of all EntityTypes.
const (
	EntityTypeVersion                      EntityType = "Version"
	EntityTypeThing                        EntityType = "Thing"
	EntityTypeLocation                     EntityType = "Location"
	EntityTypeHistoricalLocation           EntityType = "HistoricalLocation"
	EntityTypeDatastream                   EntityType = "Datastream"
	EntityTypeSensor                       EntityType = "Sensor"
	EntityTypeObservedProperty             EntityType = "ObservedProperty"
	EntityTypeObservation                  EntityType = "Observation"
	EntityTypeFeatureOfInterest            EntityType = "FeatureOfInterest"
	EntityTypeThingToLocation              EntityType = "ThingToLocation"
	EntityTypeLocationToHistoricalLocation EntityType = "LocationToHistoricalLocation"
	EntityTypeUnknown                      EntityType = "Unknown"
)

// EntityTypeList is a list for all known entity types
var EntityTypeList = []EntityType{EntityTypeThing,
	EntityTypeLocation, EntityTypeHistoricalLocation,
	EntityTypeDatastream, EntityTypeSensor,
	EntityTypeObservedProperty, EntityTypeObservation,
	EntityTypeFeatureOfInterest, EntityTypeUnknown,
}

// StringEntityMap is a map of strings that map a string to an EntityType
var StringEntityMap = map[string]EntityType{
	"thing": EntityTypeThing, "things": EntityTypeThing,
	"location": EntityTypeLocation, "locations": EntityTypeLocation,
	"historicallocation": EntityTypeHistoricalLocation, "historicallocations": EntityTypeHistoricalLocation,
	"datastream": EntityTypeDatastream, "datastreams": EntityTypeDatastream,
	"sensor": EntityTypeSensor, "sensors": EntityTypeSensor,
	"observedproperty": EntityTypeObservedProperty, "observedproperties": EntityTypeObservedProperty,
	"observation": EntityTypeObservation, "observations": EntityTypeObservation,
	"featureofinterest": EntityTypeFeatureOfInterest, "featuresofinterest": EntityTypeFeatureOfInterest,
}

// ToString return the string representation of the EntityType.
func (e EntityType) ToString() string {
	return fmt.Sprintf("%s", e)
}

// EntityFromType returns an empty entity belonging to the type
// returns nil if type cannot be mapped to an entity
func EntityFromType(e EntityType) Entity {
	switch e {
	case EntityTypeThing:
		return &Thing{}
	case EntityTypeLocation:
		return &Location{}
	case EntityTypeHistoricalLocation:
		return &HistoricalLocation{}
	case EntityTypeDatastream:
		return &Datastream{}
	case EntityTypeSensor:
		return &Sensor{}
	case EntityTypeObservedProperty:
		return &ObservedProperty{}
	case EntityTypeObservation:
		return &Observation{}
	case EntityTypeFeatureOfInterest:
		return &FeatureOfInterest{}
	}

	return nil
}

// EntityFromString returns an empty entity based on a string, returns error
// if string cannot be mapped to an entity
func EntityFromString(e string) (Entity, error) {
	et, err := EntityTypeFromString(e)
	if err != nil {
		return nil, err
	}
	return EntityFromType(et), nil
}

// EntityTypeFromString returns the EntityType for a given string
// function is case-insensitive
func EntityTypeFromString(e string) (EntityType, error) {
	val, ok := StringEntityMap[strings.ToLower(e)]
	if !ok {
		return EntityTypeUnknown, fmt.Errorf("Unknown entity %s", e)
	}

	return val, nil
}

// EntityLink holds the name and type of a SensorThings entity link.
type EntityLink string

// List of all EntityLinks.
const (
	EntityLinkThings              EntityLink = "Things"
	EntityLinkLocations           EntityLink = "Locations"
	EntityLinkHistoricalLocations EntityLink = "HistoricalLocations"
	EntityLinkDatastreams         EntityLink = "Datastreams"
	EntityLinkSensors             EntityLink = "Sensors"
	EntityLinkObservedProperties  EntityLink = "ObservedProperties"
	EntityLinkObservations        EntityLink = "Observations"
	EntityLinkFeatureOfInterests  EntityLink = "FeatureOfInterest"
)

// BaseEntity is the entry point for an entity
type BaseEntity struct {
	ID      interface{} `json:"@iot.id,omitempty"`
	NavSelf string      `json:"@iot.selfLink,omitempty"`
}

// ParseEntity defined to implement Entity
func (b *BaseEntity) ParseEntity(data []byte) error {
	return nil
}

// ContainsMandatoryParams defined to implement Entity
func (b *BaseEntity) ContainsMandatoryParams() (bool, []error) {
	return false, nil
}

// SetLinks defined to implement Entity
func (b *BaseEntity) SetLinks(externalURL string) error {
	return nil
}

// GetID return the ID of the entity
func (b *BaseEntity) GetID() interface{} {
	return b.ID
}

// GetEntityType defined to implement Entity
func (b *BaseEntity) GetEntityType() EntityType {
	return EntityTypeUnknown
}

// GetPropertyNames returns the available properties of an entity
func (b *BaseEntity) GetPropertyNames() []string {
	return nil
}

// GetSelfLink returns the self link of the entity
func (b *BaseEntity) GetSelfLink() string {
	return b.NavSelf
}

// GetSupportedEncoding defined to implement Entity
func (b *BaseEntity) GetSupportedEncoding() map[int]EncodingType {
	return nil
}

// SetID sets a newID on the entity
func (b *BaseEntity) SetID(newID interface{}) {
	b.ID = newID
}

// ToString return the string representation of the EntityLink.
func (e EntityLink) ToString() string {
	return fmt.Sprintf("%s", e)
}

// Entity is the base interface for all SensorThings entities.
type Entity interface {
	ParseEntity(data []byte) error
	ContainsMandatoryParams() (bool, []error)
	GetID() interface{}
	SetID(newID interface{})
	SetAllLinks(externalURL string)
	SetSelfLink(externalURL string)
	SetLinks(externalURL string)
	GetSelfLink() string
	GetEntityType() EntityType
	GetPropertyNames() []string
	GetSupportedEncoding() map[int]EncodingType
}

// CheckMandatoryParam checks if the given parameter is nil, if true then an ApiError will be added to the
// given list of errors.
func CheckMandatoryParam(errorList *[]error, param interface{}, entityType EntityType, paramName string) {
	isNil := false
	if param != nil {
		switch t := param.(type) {
		case string:
			if len(t) == 0 {
				isNil = true
			}
		case *string:
			if t != nil {
				t1 := *t
				if len(t1) == 0 {
					isNil = true
				}
				break
			} else {
				isNil = true
			}
		case map[string]string:
			if len(t) == 0 {
				isNil = true
			}
		case *Thing:
			var contains bool
			if t != nil {
				contains, _ = t.ContainsMandatoryParams()
			}

			if t == nil || (t.ID == nil && !contains) {
				isNil = true
			}
		case *Sensor:
			var contains bool
			if t != nil {
				contains, _ = t.ContainsMandatoryParams()
			}

			if t == nil || (t.ID == nil && !contains) {
				isNil = true
			}
		case *ObservedProperty:
			var contains bool
			if t != nil {
				contains, _ = t.ContainsMandatoryParams()
			}

			if t == nil || (t.ID == nil && !contains) {

				isNil = true
			}
		case *Datastream:
			var contains bool
			if t != nil {
				contains, _ = t.ContainsMandatoryParams()
			}

			if t == nil || (t.ID == nil && !contains) {

				isNil = true
			}
		}
	} else {
		isNil = true
	}

	err := *errorList
	if isNil {
		*errorList = append(err, gostErrors.NewBadRequestError(fmt.Errorf("Missing mandatory parameter: %s.%s", entityType, paramName)))
	}
}

// CheckEncodingSupported returns true of the Location entity supports the given encoding type
func CheckEncodingSupported(encodingType string) (bool, error) {
	_, err := CreateEncodingType(encodingType)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CreateEntitySelfLink formats the given parameters into an external navigationlink to the entity
// for example: http://example.org/OGCSensorThings/v1.0/Things(27815)
func CreateEntitySelfLink(externalURI string, entityLink string, id interface{}) string {
	if id != nil {
		entityLink = fmt.Sprintf("%s(%v)", entityLink, id)
	}

	return fmt.Sprintf("%s/v1.0/%s", externalURI, entityLink)
}

// CreateEntityLink formats the given parameters into a relative navigationlink path
// for example: http://example.org/OGCSensorThings/v1.0/Things(27815)/Datastreams
func CreateEntityLink(isNil bool, externalURI string, entityType1 string, entityType2 string, id interface{}) string {
	if !isNil {
		return ""
	}

	if id != nil {
		entityType1 = fmt.Sprintf("%s(%v)", entityType1, id)
	}

	return fmt.Sprintf("%s/v1.0/%s/%s", externalURI, entityType1, entityType2)
}
