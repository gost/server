package entities

import (
	"fmt"
	gostErrors "github.com/geodan/gost/errors"
)

// EntityType holds the name and type of a SensorThings entity.
type EntityType string

// List of all EntityTypes.
const (
	EntityTypeThing              EntityType = "Thing"
	EntityTypeLocation           EntityType = "Location"
	EntityTypeHistoricalLocation EntityType = "HistoricalLocation"
	EntityTypeDatastream         EntityType = "Datastream"
	EntityTypeSensor             EntityType = "Sensor"
	EntityTypeObservedProperty   EntityType = "ObservedProperty"
	EntityTypeObservation        EntityType = "Observation"
	EntityTypeFeatureOfInterest  EntityType = "FeatureOfInterest"
)

// ToString return the string representation of the EntityType.
func (e EntityType) ToString() string {
	return fmt.Sprintf("%s", e)
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
	EntityLinkObservedPropertys   EntityLink = "ObservedProperties"
	EntityLinkObservations        EntityLink = "Observations"
	EntityLinkFeatureOfInterests  EntityLink = "FeatureOfInterests"
)

// ToString return the string representation of the EntityLink.
func (e EntityLink) ToString() string {
	return fmt.Sprintf("%s", e)
}

// Entity is the base interface for all SensorThings entities.
type Entity interface {
	ParseEntity(data []byte) error
	ContainsMandatoryParams() (bool, []error)
	SetLinks(externalURL string)
	GetEntityType() EntityType
}

// CheckMandatoryParam checks if the given parameter is nil, if true then an ApiError will be added to the
// given list of errors.
func CheckMandatoryParam(errorList *[]error, param interface{}, entityType EntityType, paramName string) {
	isNil := false
	switch t := param.(type) {
	case string:
		if len(t) == 0 {
			isNil = true
		}
		break
	case map[string]string:
		if len(t) == 0 {
			isNil = true
		}
		break
	case *Thing:
		if t == nil || len(t.ID) == 0 {
			isNil = true
		}
		break
	case *Sensor:
		if t == nil || len(t.ID) == 0 {
			isNil = true
		}
		break
	case *ObservedProperty:
		if t == nil || len(t.ID) == 0 {
			isNil = true
		}
		break
	}

	err := *errorList
	if isNil {
		*errorList = append(err, gostErrors.NewBadRequestError(fmt.Errorf("Missing mandatory parameter: %s.%s", entityType, paramName)))
	}
}

// CreateEntitySelfLink formats the given parameters into an external navigationlink to the entity
// for example: http://example.org/OGCSensorThings/v1.0/Things(27815)
func CreateEntitySelfLink(externalURI string, entityLink string, id string) string {
	if len(id) != 0 {
		entityLink = fmt.Sprintf("%s(%s)", entityLink, id)
	}

	return fmt.Sprintf("%s/v1.0/%s", externalURI, entityLink)
}

// CreateEntityLink formats the given parameters into a relative navigationlink path
// for example: ../Things(27815)/Datastreams
func CreateEntityLink(isNil bool, entityType1 string, entityType2 string, id string) string {
	if !isNil {
		return ""
	}

	if len(id) != 0 {
		entityType1 = fmt.Sprintf("%s(%s)", entityType1, id)
	}

	return fmt.Sprintf("../%s/%s", entityType1, entityType2)
}
