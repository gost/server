package http

import (
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
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

func TestEndPointSort(t *testing.T) {
	// arrange
	restep := rest.Endpoint{}
	restep.Name = "ep1"
	// error in next line, why? About models.EndpointOperation
	//op1 := models.EndpointOperation{models.HTTPOperationGet, "/v1.0/observedproperties",HandleTest}
	//httpep1 := &Endpoint{Endpoint: restep, Operation: op1}
	httpep1 := &Endpoint{}
	httpep2 := &Endpoint{}

	eps := Endpoints{}

	eps = append(eps, httpep1)
	eps = append(eps, httpep2)

	// act
	eps.Swap(0, 1)

	// todo: add some check that the swap is ok

	// assert
	assert.True(t, len(eps) == 2, "Number of Endpoints should be 2")
}
