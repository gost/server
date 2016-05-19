package odata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderByWithEmptyOrderShouldReturnError(t *testing.T) {
	//arrange
	orderby := QueryOrderBy{}

	//act
	err := orderby.Parse("")

	//assert
	assert.NotNil(t, err, "Empty QueryOrderBy should give errors")
}

func TestOrderByParsingShouldGivePropertyAndSuffix(t *testing.T) {
	//arrange
	orderby := QueryOrderBy{}

	//act
	err := orderby.Parse("hallo asc")

	//assert
	assert.Nil(t, err, "Empty QueryOrderBy should give errors")
	assert.Equal(t, orderby.property, "hallo", "orderby.property should be correct")
	assert.Equal(t, orderby.suffix, "asc", "orderby.suffix should be correct")
}

// todo: add test for IQueryOrderBy.IsValid when implementing

func TestOrderByNotShouldReturnNotNil(t *testing.T) {
	//arrange
	orderby := QueryOrderBy{}

	//act
	res := orderby.IsNil()

	//assert
	assert.False(t, res, "OrderBy should not be nil")
}
