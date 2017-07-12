package http

import (
	"github.com/geodan/gost/sensorthings/models"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestIsDynamic(t *testing.T) {
	// arrange
	urlDynamic := "http://www.{}.nl"
	urlNotDynamic := "http://www.nu.nl"

	// act
	resultNotDynamic := isDynamic(urlNotDynamic)
	resultDynamic := isDynamic(urlDynamic)

	// assert
	assert.False(t, resultNotDynamic)
	assert.True(t, resultDynamic)
}

func TestEndPointSortDuplicate(t *testing.T) {
	// arrange
	ep1 := &Endpoint{Operation: models.EndpointOperation{Path: "ep1", OperationType: models.HTTPOperationGet}}
	ep2 := &Endpoint{Operation: models.EndpointOperation{Path: "ep1", OperationType: models.HTTPOperationGet}}
	eps := Endpoints{ep1, ep2}

	// assert
	assert.Panics(t, func() { sort.Sort(eps) })
}

func TestEndPointSort(t *testing.T) {
	// arrange
	ep1 := &Endpoint{}
	ep1.Operation.Path = "ep1"
	ep1.Operation.OperationType = models.HTTPOperationGet
	ep2 := &Endpoint{}
	ep2.Operation.Path = "ep2"
	ep2.Operation.OperationType = models.HTTPOperationPost
	ep3 := &Endpoint{}
	ep3.Operation.Path = "{c:.*}ep3"
	ep3.Operation.OperationType = models.HTTPOperationGet
	ep4 := &Endpoint{}
	ep4.Operation.Path = "ep4{c:.*}"
	ep4.Operation.OperationType = models.HTTPOperationGet
	ep5 := &Endpoint{}
	ep5.Operation.Path = "{c:.*}ep5{test}"
	ep5.Operation.OperationType = models.HTTPOperationGet
	ep6 := &Endpoint{}
	ep6.Operation.Path = "ep6{test}"
	ep6.Operation.OperationType = models.HTTPOperationGet
	ep7 := &Endpoint{}
	ep7.Operation.Path = "ep7"
	ep7.Operation.OperationType = models.HTTPOperationGet

	eps := Endpoints{ep1, ep2, ep3, ep4, ep5, ep6, ep7}

	// act
	sort.Sort(eps)

	// assert
	assert.True(t, len(eps) == 7, "Number of Endpoints should be 7")
	// post becomes first after sorting
	assert.True(t, eps[0].Operation.Path == "ep2")
	assert.True(t, eps[1].Operation.Path == "ep1")
	assert.True(t, eps[2].Operation.Path == "ep7")
	assert.True(t, eps[3].Operation.Path == "ep6{test}")
	assert.True(t, eps[4].Operation.Path == "{c:.*}ep3")
	assert.True(t, eps[5].Operation.Path == "ep4{c:.*}")
	assert.True(t, eps[6].Operation.Path == "{c:.*}ep5{test}")
}

func TestEndPointSortDynamic(t *testing.T) {
	// arrange
	httpep1 := &Endpoint{}
	httpep1.Operation.Path = "ep1{}"
	httpep1.Operation.OperationType = models.HTTPOperationGet
	httpep2 := &Endpoint{}
	httpep2.Operation.Path = "ep2{}longer"
	httpep2.Operation.OperationType = models.HTTPOperationPost

	eps := Endpoints{httpep1, httpep2}

	// act
	sort.Sort(eps)

	// assert
	assert.True(t, len(eps) == 2, "Number of Endpoints should be 2")
	// when both urls are dynamic, the longer path comes first
	assert.True(t, eps[0].Operation.Path == "ep2{}longer")
	assert.True(t, eps[1].Operation.Path == "ep1{}")
}

func TestEndPointSortlength(t *testing.T) {
	// arrange
	httpep1 := &Endpoint{}
	httpep1.Operation.Path = "ep1"
	httpep1.Operation.OperationType = models.HTTPOperationGet
	httpep2 := &Endpoint{}
	httpep2.Operation.Path = "ep2longer"
	httpep2.Operation.OperationType = models.HTTPOperationPost

	eps := Endpoints{httpep1, httpep2}

	// act
	sort.Sort(eps)

	// assert
	assert.True(t, len(eps) == 2, "Number of Endpoints should be 2")
	// when both urls are dynamic, the longer path comes first
	assert.True(t, eps[0].Operation.Path == "ep2longer")
	assert.True(t, eps[1].Operation.Path == "ep1")
}

func TestEndPointNotDynamic(t *testing.T) {
	// arrange
	httpep1 := &Endpoint{}
	httpep1.Operation.Path = "ep1 {c:.*}"
	httpep1.Operation.OperationType = models.HTTPOperationGet
	httpep2 := &Endpoint{}
	httpep2.Operation.Path = "ep2longer"
	httpep2.Operation.OperationType = models.HTTPOperationGet
	eps := Endpoints{httpep1, httpep2}

	// act
	sort.Sort(eps)

	// assert
	assert.True(t, eps[0].Operation.Path == "ep2longer")
	assert.True(t, eps[1].Operation.Path == "ep1 {c:.*}")

}
