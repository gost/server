package sensorthings

import "fmt"

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
