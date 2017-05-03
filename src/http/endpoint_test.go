package http

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestEndPointLength(t *testing.T) {
	// arrange
	ep1 := &Endpoint{}
	ep2:= &Endpoint{}
	eps := Endpoints{}
	eps = append(eps,ep1)
	eps = append(eps,ep2)

	// act
	l := eps.Len()

	// assert
	assert.True(t, l ==2, "Number of Endpoints should be 2")
}


func TestEndPointSort(t *testing.T) {
	// arrange
	ep1 := &Endpoint{}
	ep2:= &Endpoint{}
	eps := Endpoints{}
	eps = append(eps,ep1)
	eps = append(eps,ep2)

	// act
	eps.Swap(0,1)

	// todo: add some check that the swap is ok

	// assert
	assert.True(t, len(eps)==2, "Number of Endpoints should be 2")
}