package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModels(t *testing.T) {
	// todo add models tests...
	// arrange
	var a = 1
	var b = 2
	er := ErrorResponse{}
	// act
	var res = a + b
	// assert
	assert.Equal(t, 3, res, "computer error again")
	assert.NotNil(t, er)
}
