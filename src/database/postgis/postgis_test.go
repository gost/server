package postgis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestToIntIDForString(t *testing.T){
	// arrange
	fid := "4"

	// act
	intID, err  := ToIntID(fid)

	// assert
	assert.True(t,intID==4)
	assert.True(t,err)
}

func TestToIntIDForFloat(t *testing.T){
	// arrange
	fid := 6.4

	// act
	intID, err  := ToIntID(fid)

	// assert
	assert.True(t,intID==6)
	assert.True(t,err)
}

func TestContainsToLower(t *testing.T){
	// arrange
	array := []string{"Hallo"}
	search := "HALLO"

	// act
	res := ContainsToLower(array,search)

	// assert
	assert.True(t,res)
}

func TestContainsNotToLower(t *testing.T){
	// arrange
	array := []string{"Halllo"}
	search := "HALLO"

	// act
	res := ContainsToLower(array,search)

	// assert
	assert.False(t,res)
}
