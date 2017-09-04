package postgis

import (
	"net/url"
	"testing"

	entities "github.com/gost/core"
	"github.com/gost/godata"
	"github.com/gost/server/sensorthings/odata"
	"github.com/stretchr/testify/assert"
)

func TestCreateQueryBuilder(t *testing.T) {
	// act
	qb := CreateQueryBuilder("v1.0", 1)
	// assert
	assert.NotNil(t, qb)
}

func TestCreateFilter(t *testing.T) {
	qb := CreateQueryBuilder("v1.0", 1)
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenOpenParen}}, false) == "")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenCloseParen}}, false) == "")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenWhitespace}}, false) == "")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenColon}}, false) == "")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenComma}}, false) == "")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenOp}}, false) == "")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenFunc}}, false) == "")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenLambda, Value: "ho"}}, false) == "ho")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenNull, Value: "ho"}}, false) == "ho")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenIt, Value: "ho"}}, false) == "ho")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenRoot, Value: "ho"}}, false) == "ho")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenFloat, Value: "ho"}}, false) == "ho")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenInteger, Value: "ho"}}, false) == "ho")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenString, Value: "ho"}}, false) == "ho")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenDate, Value: "ho"}}, false) == "ho")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenTime, Value: "ho"}}, false) == "ho")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenDateTime, Value: "ho"}}, false) == "ho")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenBoolean, Value: "ho"}}, false) == "ho")
	assert.True(t, qb.createFilter(entities.EntityTypeThing, &godata.ParseNode{Token: &godata.Token{Type: godata.FilterTokenLiteral, Value: "ho"}}, false) == "ho")
}

func TestPrepareFilterRight(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)

	ti, ti2 := qb.prepareFilter(entities.EntityTypeDatastream, "ho", "ho", "ha", "ha")
	assert.True(t, ti == "ho" && ti2 == "ha")

	ti, ti2 = qb.prepareFilter(entities.EntityTypeDatastream, "encodingtype", "encodingtype", "application/vnd.geo+json", "application/vnd.geo+json")
	assert.True(t, ti == "encodingtype" && ti2 == "1")

	ti, ti2 = qb.prepareFilter(entities.EntityTypeDatastream, "observationtype", "observationtype", "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CategoryObservation", "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CategoryObservation")
	assert.True(t, ti == "observationtype" && ti2 == "1")

	ti, ti2 = qb.prepareFilter(entities.EntityTypeDatastream, "resulttime", "resulttime", "2006-01-02T15:04:05.000Z", "2006-01-02T15:04:05.000Z")
	assert.True(t, ti == "resulttime" && ti2 == "'2006-01-02T15:04:05.000Z'")
}

func TestRemoveSchema(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)

	// act
	res := qb.removeSchema("v2.hallo")
	// assert
	assert.True(t, res == "hallo")
}

func TestGetOffset(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	qo := &odata.QueryOptions{}
	qo.Skip, _ = godata.ParseSkipString("2")

	// act
	offset := qb.getOffset(qo)

	// assert
	assert.True(t, offset == "2")
}

func TestRemoveSchemaWithoutSchema(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)

	// act
	res := qb.removeSchema("hallo")
	// assert
	assert.True(t, res == "hallo")
}

func TestGetLimit(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	qo := &odata.QueryOptions{}

	// act
	res := qb.getLimit(qo)
	// assert
	assert.True(t, res == "1")
}

func TestGetOrderByWithNilOptions(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	ds := &entities.Datastream{}
	// act
	res := qb.getOrderBy(ds.GetEntityType(), nil)

	// assert
	assert.NotNil(t, res)
	assert.True(t, res == "datastream.id DESC")
}

func TestCreateJoin(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	thing := &entities.Thing{}
	location := &entities.Location{}

	join := qb.createJoin(thing, location, 1, false, false, nil, nil, "")
	assert.True(t, join == "INNER JOIN LATERAL (SELECT location.id AS location_id FROM v1.0.location INNER JOIN v1.0.thing_to_location ON thing_to_location.location_id = location.id AND thing_to_location.thing_id = thing.id WHERE location.id = 1) AS location on true ")
}

func TestCreateJoinWithExpand(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	thing := &entities.Thing{}
	location := &entities.Location{}

	join := qb.createJoin(thing, location, 1, true, false, nil, nil, "")
	assert.True(t, join == "LEFT JOIN LATERAL (SELECT location.id AS location_id, location.name AS location_name, location.description AS location_description, location.encodingtype AS location_encodingtype, public.ST_AsGeoJSON(location.location) AS location_location FROM v1.0.location INNER JOIN v1.0.thing_to_location ON thing_to_location.location_id = location.id AND thing_to_location.thing_id = thing.id  ORDER BY location.id DESC LIMIT 1 OFFSET 0) AS location on true ")
}
func TestCreateCountQuery(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	expected := "SELECT COUNT(*) FROM v1.0.datastream INNER JOIN LATERAL (SELECT thing.id AS thing_id FROM v1.0.thing WHERE thing.id = datastream.thing_id AND thing.id = 1) AS thing on true  WHERE thing.thing_id = 1 AND  datastream.name = 'Milk' AND Price < 2.55"
	qo := &odata.QueryOptions{}
	input := "Name eq 'Milk' and Price lt 2.55"
	filter, _ := godata.ParseFilterString(input)
	qo.Filter = filter

	res := qb.CreateCountQuery(&entities.Datastream{}, &entities.Thing{}, 1, qo)

	// assert
	assert.NotNil(t, res)
	assert.True(t, expected == res)
}

func TestGetOrderByWithQueryOptions(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	qo := &odata.QueryOptions{}
	qob, _ := godata.ParseOrderByString("id asc,name desc")
	qo.OrderBy = qob
	ds := &entities.Datastream{}

	// act
	res := qb.getOrderBy(ds.GetEntityType(), qo)

	// assert
	assert.NotNil(t, res)
	assert.True(t, res == "datastream.id asc, datastream.name desc")
}

func TestGetLimitWithQueryTop(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	qo := &odata.QueryOptions{}
	top, _ := godata.ParseTopString("2")
	qo.Top = top

	// act
	res := qb.getLimit(qo)
	// assert
	assert.True(t, res == "2")
}

func TestOdataLogicalOperatorToPostgreSQL(t *testing.T) {
	qb := CreateQueryBuilder("v1.0", 1)
	assert.True(t, qb.odataLogicalOperatorToPostgreSQL("and") == "AND")
	assert.True(t, qb.odataLogicalOperatorToPostgreSQL("or") == "OR")
	assert.True(t, qb.odataLogicalOperatorToPostgreSQL("not") == "NOT")
	assert.True(t, qb.odataLogicalOperatorToPostgreSQL("has") == "HAS")
	assert.True(t, qb.odataLogicalOperatorToPostgreSQL("ne") == "!=")
	assert.True(t, qb.odataLogicalOperatorToPostgreSQL("gt") == ">")
	assert.True(t, qb.odataLogicalOperatorToPostgreSQL("ge") == ">=")
	assert.True(t, qb.odataLogicalOperatorToPostgreSQL("lt") == "<")
	assert.True(t, qb.odataLogicalOperatorToPostgreSQL("le") == "<=")
	assert.True(t, qb.odataLogicalOperatorToPostgreSQL("ho") == "")
}

func TestCreateCountQueryWithoutId(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	expected := "SELECT COUNT(*) FROM v1.0.datastream INNER JOIN LATERAL (SELECT thing.id AS thing_id FROM v1.0.thing WHERE thing.id = datastream.thing_id ) AS thing on true  WHERE datastream.name = 'Milk' AND Price < 2.55"
	qo := &odata.QueryOptions{}
	input := "Name eq 'Milk' and Price lt 2.55"
	filter, _ := godata.ParseFilterString(input)
	qo.Filter = filter

	res := qb.CreateCountQuery(&entities.Datastream{}, &entities.Thing{}, nil, qo)

	// assert
	assert.NotNil(t, res)
	assert.True(t, expected == res)
}

func TestCreateQuery(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	expected := "SELECT A_datastream.datastream_id AS A_datastream_id, A_datastream.datastream_name AS A_datastream_name, A_datastream.datastream_description AS A_datastream_description, A_datastream.datastream_unitofmeasurement AS A_datastream_unitofmeasurement, A_datastream.datastream_observationtype AS A_datastream_observationtype, A_datastream.datastream_observedarea AS A_datastream_observedarea, A_datastream.datastream_phenomenontime AS A_datastream_phenomenontime, A_datastream.datastream_resulttime AS A_datastream_resulttime FROM (SELECT datastream.thing_id AS datastream_thing_id, datastream.observedproperty_id AS datastream_observedproperty_id, datastream.sensor_id AS datastream_sensor_id, datastream.id AS datastream_id, datastream.name AS datastream_name, datastream.description AS datastream_description, datastream.unitofmeasurement AS datastream_unitofmeasurement, datastream.observationtype AS datastream_observationtype, public.ST_AsGeoJSON(datastream.observedarea) AS datastream_observedarea, datastream.phenomenontime AS datastream_phenomenontime, datastream.resulttime AS datastream_resulttime FROM v1.0.datastream ORDER BY datastream.id DESC ) AS A_datastream INNER JOIN LATERAL (SELECT thing.id AS thing_id FROM v1.0.thing WHERE thing.id = A_datastream.datastream_thing_id AND thing.id = 0) AS thing on true   OFFSET 0"

	// act
	query, _ := qb.CreateQuery(&entities.Datastream{}, &entities.Thing{}, 0, nil)

	// assert
	assert.NotNil(t, query)
	assert.True(t, expected == query)
}

func TestConstructQueryParseInfo(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	expandItem1 := &godata.ExpandItem{}
	token := &godata.Token{}
	token.Value = "thing"
	tokens := []*godata.Token{token}
	expandItem1.Path = tokens
	expandItems := []*godata.ExpandItem{expandItem1}
	qpi := &QueryParseInfo{}

	// act
	qb.constructQueryParseInfo(expandItems, qpi)

	// assert
}

func TestSortQueryOptionsScenario1(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1", 1)
	uri, _ := url.Parse("localhost/v1.0/Observations?$filter=Datastream/id eq 1")
	qo, _ := odata.ParseURLQuery(uri.Query())

	// act
	qb.sortQueryOptions(qo)

	// assert
	assert.Equal(t, "Datastream", qo.Expand.ExpandItems[0].Path[0].Value)
	assert.Equal(t, "eq", qo.Expand.ExpandItems[0].Filter.Tree.Token.Value)
	assert.Equal(t, "id", qo.Expand.ExpandItems[0].Filter.Tree.Children[0].Token.Value)
	assert.Equal(t, "1", qo.Expand.ExpandItems[0].Filter.Tree.Children[1].Token.Value)
}

func TestSortQueryOptionsScenario2(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1", 1)
	uri, _ := url.Parse("localhost/v1.0/Things?$filter=Datastreams/Observations/result ge 20")
	qo, _ := odata.ParseURLQuery(uri.Query())

	// act
	qb.sortQueryOptions(qo)

	// assert
	assert.Equal(t, "Datastreams", qo.Expand.ExpandItems[0].Path[0].Value)
	assert.Equal(t, "Observations", qo.Expand.ExpandItems[0].Path[1].Value)
	assert.Equal(t, "ge", qo.Expand.ExpandItems[0].Filter.Tree.Token.Value)
	assert.Equal(t, "result", qo.Expand.ExpandItems[0].Filter.Tree.Children[0].Token.Value)
	assert.Equal(t, "20", qo.Expand.ExpandItems[0].Filter.Tree.Children[1].Token.Value)
}

func TestSortQueryOptionsScenario3(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1", 1)
	uri, _ := url.Parse("localhost/v1.0/Things?$expand=Locations,Datastreams/Observations&$filter=length(Datastreams/name) eq 10 and Datastreams/Observations/result ge 20")
	qo, _ := odata.ParseURLQuery(uri.Query())

	// act
	qb.sortQueryOptions(qo)

	// assert
	assert.Equal(t, 0, len(qo.Filter.Tree.Children))

	assert.Equal(t, "Locations", qo.Expand.ExpandItems[0].Path[0].Value)

	assert.Equal(t, "Datastreams", qo.Expand.ExpandItems[1].Path[0].Value)
	assert.Equal(t, "Observations", qo.Expand.ExpandItems[1].Path[1].Value)
	assert.Equal(t, "ge", qo.Expand.ExpandItems[1].Filter.Tree.Token.Value)
	assert.Equal(t, "result", qo.Expand.ExpandItems[1].Filter.Tree.Children[0].Token.Value)
	assert.Equal(t, "20", qo.Expand.ExpandItems[1].Filter.Tree.Children[1].Token.Value)

	assert.Equal(t, "Datastreams", qo.Expand.ExpandItems[2].Path[0].Value)
	assert.Equal(t, "eq", qo.Expand.ExpandItems[2].Filter.Tree.Token.Value)
	assert.Equal(t, "length", qo.Expand.ExpandItems[2].Filter.Tree.Children[0].Token.Value)
	assert.Equal(t, "name", qo.Expand.ExpandItems[2].Filter.Tree.Children[0].Children[0].Token.Value)
	assert.Equal(t, "10", qo.Expand.ExpandItems[2].Filter.Tree.Children[1].Token.Value)
}

func TestSortQueryOptionsScenario4(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1", 1)
	uri, _ := url.Parse("localhost/v1.0/Observations?$filter=Datastream/id eq 1 and Datastream/name eq test or length(Datastream/name) eq 10")
	qo, _ := odata.ParseURLQuery(uri.Query())

	// act
	qb.sortQueryOptions(qo)

	// assert
	assert.Equal(t, "Datastream", qo.Expand.ExpandItems[0].Path[0].Value)
	assert.Equal(t, "or", qo.Expand.ExpandItems[0].Filter.Tree.Token.Value)
	assert.Equal(t, "and", qo.Expand.ExpandItems[0].Filter.Tree.Children[0].Token.Value)

	assert.Equal(t, "eq", qo.Expand.ExpandItems[0].Filter.Tree.Children[0].Children[0].Token.Value)
	assert.Equal(t, "id", qo.Expand.ExpandItems[0].Filter.Tree.Children[0].Children[0].Children[0].Token.Value)
	assert.Equal(t, "1", qo.Expand.ExpandItems[0].Filter.Tree.Children[0].Children[0].Children[1].Token.Value)

	assert.Equal(t, "eq", qo.Expand.ExpandItems[0].Filter.Tree.Children[0].Children[1].Token.Value)
	assert.Equal(t, "name", qo.Expand.ExpandItems[0].Filter.Tree.Children[0].Children[1].Children[0].Token.Value)
	assert.Equal(t, "test", qo.Expand.ExpandItems[0].Filter.Tree.Children[0].Children[1].Children[1].Token.Value)

	assert.Equal(t, "eq", qo.Expand.ExpandItems[0].Filter.Tree.Children[1].Token.Value)
	assert.Equal(t, "length", qo.Expand.ExpandItems[0].Filter.Tree.Children[1].Children[0].Token.Value)
	assert.Equal(t, "name", qo.Expand.ExpandItems[0].Filter.Tree.Children[1].Children[0].Children[0].Token.Value)
	assert.Equal(t, "10", qo.Expand.ExpandItems[0].Filter.Tree.Children[1].Children[1].Token.Value)
}
