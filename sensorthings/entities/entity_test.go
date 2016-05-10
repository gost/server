package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntities(t *testing.T) {
	// todo add entities tests...
	// arrange
	var a = 1
	var b = 2
	// act
	var res = a + b
	// assert
	assert.Equal(t, 3, res, "computer error again")
}
