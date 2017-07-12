package odata

import (
	"fmt"
	"github.com/gost/godata"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestExpandParametersSupported(t *testing.T) {
	// arrange
	qo := &QueryOptions{}
	SupportedExpandParameters = map[string][]string{"things": {"locations", "datastreams"}}

	// assert
	assert.Equal(t, true, qo.ExpandParametersSupported("things", "locations"))
	assert.Equal(t, false, qo.ExpandParametersSupported("things", "featuresofinterest"))
	assert.Equal(t, false, qo.ExpandParametersSupported("bla", ""))
}

func TestSelectParametersSupported(t *testing.T) {
	// arrange
	qo := &QueryOptions{}
	SupportedSelectParameters = map[string][]string{"things": {"id", "name"}}

	// assert
	assert.Equal(t, true, qo.SelectParametersSupported("things", "name"))
	assert.Equal(t, false, qo.SelectParametersSupported("things", "nonexistingparam"))
	assert.Equal(t, false, qo.SelectParametersSupported("bla", ""))
}

func TestEmptyQuery(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Nil(t, query, "parsing query from localhost/v1.0/things should return nil")
}

func TestWrongQuery(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things?$count=ok")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Nil(t, query, "parsing query from localhost/v1.0/things?$count=ok should return nil because ok cannot be parsed to bool")
}

func TestParseFilter(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things?$filter=id%20eq%201&$top=1&$skip=1&$count=true&$expand=Datastreams/Observations,Locations&$orderby=id&$select=name")

	// act
	query, err := ParseURLQuery(uri.Query())

	// assert
	assert.NotNil(t, query, fmt.Sprintf("%v", err))
	assert.NotNil(t, query.Filter)
	assert.NotNil(t, query.Count)
	assert.NotNil(t, query.Expand)
	assert.NotNil(t, query.OrderBy)
	assert.NotNil(t, query.Select)
	assert.NotNil(t, query.Top)
	assert.NotNil(t, query.Skip)
}

func TestParseFilterRef(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things?$ref=true")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Equal(t, GoDataRefQuery(true), *query.Ref)
}

func TestParseFilterValue(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things?$value=true")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Equal(t, GoDataValueQuery(true), *query.Value)
}

func TestSavingRawQuery(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things?$filter=id eq 1&$orderby=id desc&$expand=Datastreams/Observations,Locations")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Equal(t, "id eq 1", query.RawFilter)
	assert.Equal(t, "id desc", query.RawOrderBy)
	assert.Equal(t, "Datastreams/Observations,Locations", query.RawExpand)
}

func TestFilterGeography(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/Locations?$filter=geo.intersects(location,geography'LINESTRING(7.5 51.5, 7.5 53.5)')")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Equal(t, godata.FilterTokenFunc, query.Filter.Tree.Token.Type)
	assert.Equal(t, godata.FilterTokenLiteral, query.Filter.Tree.Children[0].Token.Type)
	assert.Equal(t, godata.FilterTokenGeography, query.Filter.Tree.Children[1].Token.Type)
}

func TestExpandItemToQueryOptions(t *testing.T) {
	// arrange
	ei := &godata.ExpandItem{Filter: &godata.GoDataFilterQuery{}}

	// act
	qo := ExpandItemToQueryOptions(ei)

	// assert
	assert.Equal(t, ei.Filter, qo.Filter)
}
