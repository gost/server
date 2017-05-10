package postgis

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateTableMappings(t *testing.T){
	// arrange

	// act
	mappings :=createTableMappings("v1")

	// assert
	assert.NotNil(t, mappings)
	assert.True(t,len(mappings)>0)
}
