package rest

import (
	"testing"

	"github.com/geodan/gost/src/sensorthings/odata"
	"github.com/stretchr/testify/assert"
)

func TestEndPointGetNameShouldReturnCorrectName(t *testing.T) {
	//arrange
	endpoint := Endpoint{}
	endpoint.Name = "test"

	// act
	name := endpoint.GetName()

	//assert
	assert.Equal(t, "test", name, "name shouldbe correct")
}

func TestEndPointGetQueryOptions(t *testing.T) {
	// arrange
	qo := &odata.QueryOptions{}
	endpoint := Endpoint{}

	//act
	b, _ := endpoint.AreQueryOptionsSupported(qo)

	// assert
	assert.True(t, b, "QueryOptionsSupport should be true")
}
