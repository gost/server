package sensorthings

import "fmt"

func CreateEntitySefLink(externalUri string, entityType string, id string) string{
	if(len(id) != 0){
		entityType = fmt.Sprintf("%s(%s)", entityType, id)
	}

	return fmt.Sprintf("%s/v1.0/%s", externalUri, entityType)
}

func CreateEntityLink(entityType1 string, entityType2 string, id string) string{
	if(len(id) != 0){
		entityType1 = fmt.Sprintf("%s(%s)", entityType1, id)
	}

	return fmt.Sprintf("../%s/%s", entityType1, entityType2)
}