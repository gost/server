package odata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpandParserShouldFillParams(t *testing.T) {
	//arrange
	expand := QueryExpand{}

	//act
	res := expand.Parse("fld1,fld2")

	//assert
	assert.Nil(t, res, "Parse result should be nil")
	assert.True(t, len(expand.params) > 0, "Expand params should be filled")
}

func TestExpandIsValidShouldReturnTrueWhenValid(t *testing.T) {
	//arrange
	expand := QueryExpand{}
	expand.Parse("fld1,fld2")
	var vals = []string{"fld1", "fld2"}

	//act
	bln, _ := expand.IsValid(vals, "haha")

	//assert
	assert.True(t, bln, "Expand IsValid should return correct value")
}

func TestExpandIsValidShouldReturnFalseWhenNotValid(t *testing.T) {
	//arrange
	expand := QueryExpand{}
	expand.Parse("fld1,fld2")
	var vals = []string{"hoho", "haha"}

	//act
	bln, _ := expand.IsValid(vals, "haha")

	//assert
	assert.False(t, bln, "Expand IsValid should return correct value")
}

func TestExpandisNilNotShouldReturnNotNil(t *testing.T) {
	//arrange
	expand := QueryOrderBy{}

	//act
	res := expand.IsNil()

	//assert
	assert.False(t, res, "Expand should not be nil")
}
