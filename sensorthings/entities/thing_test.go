package entities

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMissingMandatoryParameters(t *testing.T) {
	//arrange
	thing := &Thing{}

	//act
	_, err := thing.ContainsMandatoryParams()

	//assert
	assert.NotNil(t, err, "Thing mandatory param description not filled in should have returned error")
	if len(err) > 0 {
		assert.Contains(t, fmt.Sprintf("%v", err[0]), "description")
	}
}

func TestMandatoryParametersExist(t *testing.T) {
	//arrange
	thing := &Thing{Description: "test"}

	//act
	_, err := thing.ContainsMandatoryParams()

	//assert
	assert.Nil(t, err, "All mandatory params are fiulled in shoud not have return an error")
}
