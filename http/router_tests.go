package http

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestHttp starts
func TestHttp(t *testing.T) {
    // todo add http tests...
    // arrange
    var a=1
    var b=2
    // act
    var res = a+b
    // assert
	assert.Equal(t,3, res, "computer error again")
}