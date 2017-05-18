package odata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetQueryOptionTypeTop(t *testing.T) {
	//arrange
	top := QueryTop{}

	//act
	optionType := top.GetQueryOptionType()

	//assert
	assert.Equal(t, QueryOptionTop, optionType, "QueryTop should be of type QueryOptionTop")
}

func TestParseTop(t *testing.T) {
	//arrange
	top := QueryTop{}

	//act
	err := top.Parse("120")

	//assert
	assert.Nil(t, err)
	assert.Equal(t, 120, top.Limit, "QueryTop.limit should have been 120")
}

func TestParseFailTop(t *testing.T) {
	//arrange
	top := QueryTop{}

	//act
	err := top.Parse("one")

	//assert
	assert.NotNil(t, err)
}

func TestIsNilTrueTop(t *testing.T) {
	//arrange
	qo := QueryOptions{}

	//assert
	assert.Equal(t, true, qo.QueryTop.IsNil())
}

func TestIsNilFalseTop(t *testing.T) {
	//arrange
	qo := QueryOptions{}
	qo.QueryTop = &QueryTop{}

	//act
	qo.QueryTop.Parse("1")

	//assert
	assert.Equal(t, false, qo.QueryTop.IsNil())
}
