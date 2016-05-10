package odata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCount(t *testing.T) {
	// todo add count tests...
	// arrange
	var a = 1
	var b = 2
	// act
	var res = a + b
	// assert
	assert.Equal(t, 3, res, "computer error again")
}
