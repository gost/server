package entities

import (
	"errors"
	"fmt"
	gostErrors "github.com/geodan/gost/src/errors"
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
	EntityTypeUnknown            EntityType = "Unknown"
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

// BaseEntity is the entry point for an entity
type BaseEntity struct {
	ID      string `json:"@iot.id,omitempty"`
	NavSelf string `json:"@iot.selfLink,omitempty"`
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

// ToString return the string representation of the EntityLink.
func (e EntityLink) ToString() string {
	return fmt.Sprintf("%s", e)
}

// Entity is the base interface for all SensorThings entities.
type Entity interface {
	ParseEntity(data []byte) error
	ContainsMandatoryParams() (bool, []error)
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
			break
		case map[string]string:
			if len(t) == 0 {
				isNil = true
			}
			break
		case *Thing:
			var contains bool
			if t != nil {
				contains, _ = t.ContainsMandatoryParams()
			}

			if t == nil || (len(t.ID) == 0 && !contains) {
				isNil = true
			}
			break
		case *Sensor:
			var contains bool
			if t != nil {
				contains, _ = t.ContainsMandatoryParams()
			}

			if t == nil || (len(t.ID) == 0 && !contains) {

				isNil = true
			}
			break
		case *ObservedProperty:
			var contains bool
			if t != nil {
				contains, _ = t.ContainsMandatoryParams()
			}

			if t == nil || (len(t.ID) == 0 && !contains) {

				isNil = true
			}
			break
		case *Datastream:
			var contains bool
			if t != nil {
				contains, _ = t.ContainsMandatoryParams()
			}

			if t == nil || (len(t.ID) == 0 && !contains) {

				isNil = true
			}
			break
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
func CheckEncodingSupported(entity Entity, encodingType string) (bool, error) {
	notSupported := gostErrors.NewBadRequestError(errors.New("encodingType not supported"))
	encoding, err := CreateEncodingType(encodingType)
	if err != nil {
		return false, notSupported
	}

	supportedEncodings := entity.GetSupportedEncoding()
	_, ok := supportedEncodings[encoding.Code]
	if ok {
		return true, nil
	}

	return false, notSupported
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
// for example: http://example.org/OGCSensorThings/v1.0/Things(27815)/Datastreams
func CreateEntityLink(isNil bool, externalURI string, entityType1 string, entityType2 string, id string) string {
	if !isNil {
		return ""
	}

	if len(id) != 0 {
		entityType1 = fmt.Sprintf("%s(%s)", entityType1, id)
	}

	return fmt.Sprintf("%s/v1.0/%s/%s", externalURI, entityType1, entityType2)
}
