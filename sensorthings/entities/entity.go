package entities

import "fmt"

// EntityType holds the name and type of an SensorThings entity
type EntityType string

// List of all EntityTypes
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

// Entity is the base interface for all SensorThings entities
type Entity interface {
	ParseEntity(data []byte) error
	ContainsMandatoryPostParams() (bool, []error)
	SetLinks(externalURL string)
}

// CreateEntitySefLink formats the given parameters into an external navigationlink to the entity
// for example: http://example.org/OGCSensorThings/v1.0/Things(27815)
func CreateEntitySefLink(externalURI string, entityType string, id string) string {
	if len(id) != 0 {
		entityType = fmt.Sprintf("%s(%s)", entityType, id)
	}

	return fmt.Sprintf("%s/v1.0/%s", externalURI, entityType)
}

// CreateEntityLink formats the given parameters into a relative navigationlink path
// for example: ../Things(27815)/Datastreams
func CreateEntityLink(entityType1 string, entityType2 string, id string) string {
	if len(id) != 0 {
		entityType1 = fmt.Sprintf("%s(%s)", entityType1, id)
	}

	return fmt.Sprintf("../%s/%s", entityType1, entityType2)
}
