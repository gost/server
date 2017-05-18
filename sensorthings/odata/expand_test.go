package odata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpandParserShouldFillParams(t *testing.T) {
	//arrange
	expand := QueryExpand{}

	//act
	res := expand.Parse("Thing,Location")

	//assert
	assert.Nil(t, res, "Parse result should be nil")
	assert.True(t, len(expand.Operations) > 0, "Expand operations should be filled")
}

func TestExpandisNilNotShouldReturnNotNil(t *testing.T) {
	//arrange
	expand := QueryOrderBy{}

	//act
	res := expand.IsNil()

	//assert
	assert.False(t, res, "Expand should not be nil")
}
