package rest

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
    // todo add config tests...
    // arrange
    var a=1
    var b=2
    // act
    var res = a+b
    // assert
	assert.Equal(t,3, res, "computer error again")
}