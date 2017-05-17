package rest

import (
	"testing"

	"github.com/geodan/gost/src/sensorthings/odata"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func TestEndPointGetNameShouldReturnCorrectName(t *testing.T) {
	//arrange
	endpoint := Endpoint{}
	endpoint.Name = "test"
	endpoint.URL = "http://www.nu.nl"

	// act
	name := endpoint.GetName()
	output := endpoint.ShowOutputInfo()
	url := endpoint.GetURL()
	ops := endpoint.GetOperations()
	qos := endpoint.GetSupportedQueryOptions()
	expand := endpoint.GetSupportedExpandParams()
	sel := endpoint.GetSupportedSelectParams()
	// point.AreQueryOptionsSupported()

	//assert
	assert.Equal(t, "test", name, "name should be correct")
	assert.True(t, !output)
	assert.Equal(t, url, "http://www.nu.nl")
	assert.True(t, len(ops) == 0)
	assert.True(t, len(qos) == 0)
	assert.True(t, len(expand) == 0)
	assert.True(t, len(sel) == 0)

}

func TestEndPointGetQueryOptions(t *testing.T) {
	// arrange
	qo := &odata.QueryOptions{}
	endpoint := Endpoint{}

	// act
	b, _ := endpoint.AreQueryOptionsSupported(qo)

	// assert
	assert.True(t, b, "QueryOptionsSupport should be true")
}

func TestSupportsQueryOptionTypeResultFalse(t *testing.T) {
	// arrange
	qo := odata.QueryOptions{}
	endpoint := Endpoint{}

	// act
	b := endpoint.SupportsQueryOptionType(qo.QueryCount.GetQueryOptionType())

	// assert
	assert.False(t, b)
}

func TestQueryOptionSupported (t *testing.T) {
	// arrange
	epDatastream := createDatastreamsEndpoint("http:/www.nu.nl")
	qo := odata.QueryOptions{}
	qo.QueryCount = &odata.QueryCount{}
	qo.QueryExpand = &odata.QueryExpand{}

	var errorList []error

	// act
	checkQueryOptionSupported(epDatastream, qo.QueryCount, &errorList, odata.CreateQueryError(odata.QueryNotAvailable, http.StatusNotImplemented, qo.QueryTop.GetQueryOptionType().String(), epDatastream.Name))
	//checkQueryOptionSupported(epDatastream, qo.QueryExpand, &errorList, odata.CreateQueryError(odata.QueryNotAvailable, http.StatusNotImplemented, qo.QueryTop.GetQueryOptionType().String(), epDatastream.Name))

	// assert
	assert.True(t, len(errorList) == 0)
}

func TestSupportsQueryOptionTypeResultTrue(t *testing.T) {
	// arrange
	epDatatream := createDatastreamsEndpoint("http:/www.nu.nl")
	querycount := odata.QueryOptions{}.QueryExpand

	//act
	supportsQueryCount := epDatatream.SupportsQueryOptionType(querycount.GetQueryOptionType())
	// assert
	assert.True(t, supportsQueryCount)
}
