package odata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetQueryOptionTypeCount(t *testing.T) {
	//arrange
	qOptCount := QueryCount{}

	//act
	optionType := qOptCount.GetQueryOptionType()

	//assert
	assert.Equal(t, QueryOptionCount, optionType, "QueryCount should be of type QueryOptionCount")
}

func TestParseCount(t *testing.T) {
	//arrange
	qOptCount := QueryCount{}

	//act
	err := qOptCount.Parse("True")

	//assert
	assert.Nil(t, err)
	assert.Equal(t, true, qOptCount.count, "QueryCount.count should have been true")
}

func TestParseFailCount(t *testing.T) {
	//arrange
	qOptCount := QueryCount{}

	//act
	err := qOptCount.Parse("tru")

	//assert
	assert.NotNil(t, err)
	assert.Equal(t, false, qOptCount.count, "QueryCount.count should have been false")
}

func TestIsNilTrueCount(t *testing.T) {
	//arrange
	qo := QueryOptions{}

	//assert
	assert.Equal(t, true, qo.QueryCount.IsNil())
}

func TestIsNilFalseCount(t *testing.T) {
	//arrange
	qo := QueryOptions{}
	qo.QueryCount = &QueryCount{}

	//act
	qo.QueryCount.Parse("true")

	//assert
	assert.Equal(t, false, qo.QueryCount.IsNil())
}
