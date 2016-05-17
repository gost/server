package odata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetQueryOptionTypeSkip(t *testing.T) {
	//arrange
	skip := QuerySkip{}

	//act
	optionType := skip.GetQueryOptionType()

	//assert
	assert.Equal(t, QueryOptionSkip, optionType, "QuerySkip should be of type QueryOptionSkip")
}

func TestParseSkip(t *testing.T) {
	//arrange
	skip := QuerySkip{}

	//act
	err := skip.Parse("435")

	//assert
	assert.Nil(t, err)
	assert.Equal(t, 435, skip.index, "QuerySkip.index should have been 435")
}

func TestParseFailSkip(t *testing.T) {
	//arrange
	skip := QuerySkip{}

	//act
	err := skip.Parse("100.")

	//assert
	assert.NotNil(t, err)
}

func TestIsNilTrueSkip(t *testing.T) {
	//arrange
	qo := QueryOptions{}

	//assert
	assert.Equal(t, true, qo.QuerySkip.IsNil())
}

func TestIsNilFalseSkip(t *testing.T) {
	//arrange
	qo := QueryOptions{}
	qo.QuerySkip = &QuerySkip{}

	//act
	qo.QuerySkip.Parse("10")

	//assert
	assert.Equal(t, false, qo.QuerySkip.IsNil())
}
