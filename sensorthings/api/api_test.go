package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
    // todo add api tests...
    // arrange
    var a=1
    var b=2
    // act
    var res = a+b
    // assert
	assert.Equal(t,3, res, "computer error again")
}