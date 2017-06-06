package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndPointGetNameShouldReturnCorrectName(t *testing.T) {
	//arrange
	endpoint := Endpoint{}
	endpoint.Name = "test"
	endpoint.URL = "http://www.nu.nl"

	// act
	name := endpoint.GetName()
	output := endpoint.ShowOutputInfo()
	url := endpoint.GetURL()
	ops := endpoint.GetOperations()
	expand := endpoint.GetSupportedExpandParams()
	sel := endpoint.GetSupportedSelectParams()
	// point.AreQueryOptionsSupported()

	//assert
	assert.Equal(t, "test", name, "name should be correct")
	assert.True(t, !output)
	assert.Equal(t, url, "http://www.nu.nl")
	assert.True(t, len(ops) == 0)
	assert.True(t, len(expand) == 0)
	assert.True(t, len(sel) == 0)

}
