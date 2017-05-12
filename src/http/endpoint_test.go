package http

import (
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sort"
	"testing"
)

func TestEndPointLength(t *testing.T) {
	// arrange
	ep1 := &Endpoint{}
	ep2 := &Endpoint{}
	eps := Endpoints{}
	eps = append(eps, ep1)
	eps = append(eps, ep2)

	// act
	l := eps.Len()

	// assert
	assert.True(t, l == 2, "Number of Endpoints should be 2")
}

func HandleTest(w http.ResponseWriter, r *http.Request, endpoint *models.Endpoint, api *models.API) {
}

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

func TestEndPointSort(t *testing.T) {
	// arrange
	httpep1 := &Endpoint{}
	httpep1.Operation.Path = "ep1"
	httpep1.Operation.OperationType = models.HTTPOperationGet
	httpep2 := &Endpoint{}
	httpep2.Operation.Path = "ep2"
	httpep2.Operation.OperationType = models.HTTPOperationPost

	eps := Endpoints{httpep1,httpep2}

	// act
	sort.Sort(eps)

	// assert
	assert.True(t, len(eps) == 2, "Number of Endpoints should be 2")
	// post becomes first after sorting
	assert.True(t, eps[0].Operation.Path=="ep2")
	assert.True(t, eps[1].Operation.Path=="ep1")
}


func TestEndPointSortDynamic(t *testing.T) {
	// arrange
	httpep1 := &Endpoint{}
	httpep1.Operation.Path = "ep1{}"
	httpep1.Operation.OperationType = models.HTTPOperationGet
	httpep2 := &Endpoint{}
	httpep2.Operation.Path = "ep2{}longer"
	httpep2.Operation.OperationType = models.HTTPOperationPost

	eps := Endpoints{httpep1,httpep2}

	// act
	sort.Sort(eps)

	// assert
	assert.True(t, len(eps) == 2, "Number of Endpoints should be 2")
	// when both urls are dynamic, the longer path comes first
	assert.True(t, eps[0].Operation.Path=="ep2{}longer")
	assert.True(t, eps[1].Operation.Path=="ep1{}")
}

func TestEndPointSortlength(t *testing.T) {
	// arrange
	httpep1 := &Endpoint{}
	httpep1.Operation.Path = "ep1"
	httpep1.Operation.OperationType = models.HTTPOperationGet
	httpep2 := &Endpoint{}
	httpep2.Operation.Path = "ep2longer"
	httpep2.Operation.OperationType = models.HTTPOperationPost

	eps := Endpoints{httpep1,httpep2}

	// act
	sort.Sort(eps)

	// assert
	assert.True(t, len(eps) == 2, "Number of Endpoints should be 2")
	// when both urls are dynamic, the longer path comes first
	assert.True(t, eps[0].Operation.Path=="ep2longer")
	assert.True(t, eps[1].Operation.Path=="ep1")
}
