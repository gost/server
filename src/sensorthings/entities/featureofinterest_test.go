package entities


import (
"fmt"
"github.com/stretchr/testify/assert"
"testing"
)

func TestMissingMandatoryParametersFeatureOfInterest(t *testing.T) {
	//arrange
	featureofinterest := &FeatureOfInterest{}

	//act
	_, err := featureofinterest.ContainsMandatoryParams()

	//assert
	assert.NotNil(t, err, "FeatureOfInterest mandatory param description not filled in should have returned error")
	if len(err) > 0 {
		assert.Contains(t, fmt.Sprintf("%v", err[0]), "description")
	}
}


