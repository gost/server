package postgis

import (
	entities "github.com/gost/core"
	"github.com/gost/godata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTableMappings(t *testing.T) {
	// arrange

	// act
	mappings := createTableMappings("v1")

	// assert
	assert.NotNil(t, mappings)
	assert.True(t, len(mappings) > 0)
	assert.True(t, mappings["Thing"] == "v1.thing")
}

func TestQueryParseInfoInit(t *testing.T) {
	// arrange
	qpi := QueryParseInfo{}
	pqpi := &QueryParseInfo{}

	eo := &godata.ExpandItem{}

	// act
	qpi.Init(entities.EntityTypeThing, 0, pqpi, eo)
	qpi.Init(entities.EntityTypeFeatureOfInterest, 0, pqpi, eo)
	qpi.Init(entities.EntityTypeLocation, 0, pqpi, eo)
	qpi.Init(entities.EntityTypeObservation, 0, pqpi, eo)
	qpi.Init(entities.EntityTypeObservedProperty, 0, pqpi, eo)
	qpi.Init(entities.EntityTypeDatastream, 0, pqpi, eo)
	qpi.Init(entities.EntityTypeHistoricalLocation, 0, pqpi, eo)
	qpi.Init(entities.EntityTypeSensor, 0, pqpi, eo)

	// assert
	assert.True(t, true)
}

func TestGetQueryParseInfoByQueryIndex(t *testing.T) {
	// arrange
	qpi := QueryParseInfo{}

	// act
	qpiReturn := qpi.GetQueryParseInfoByQueryIndex(0)

	// assert
	assert.NotNil(t, qpiReturn)
	assert.NotNil(t, qpiReturn.QueryIndex == 0)

}

func TestGetQueryParseInfoByQueryIndexWithSubEntities(t *testing.T) {
	// arrange
	qpi := QueryParseInfo{}
	subqpi := &QueryParseInfo{}
	subqpi.QueryIndex = 1
	qpi.SubEntities = []*QueryParseInfo{subqpi}

	// act
	qpiReturn := qpi.GetQueryParseInfoByQueryIndex(1)

	// assert
	assert.True(t, qpiReturn.QueryIndex == 1)
}

func TestGetQueryParseInfoByQueryIndexNotFound(t *testing.T) {
	// arrange
	qpi := QueryParseInfo{}
	subqpi := &QueryParseInfo{}
	subqpi.QueryIndex = 1
	qpi.SubEntities = []*QueryParseInfo{subqpi}

	// act
	qpiReturn := qpi.GetQueryParseInfoByQueryIndex(2)

	// assert
	assert.Nil(t, qpiReturn)
}

func TestGetNextQueryIndex(t *testing.T) {
	// arrange
	qpi := QueryParseInfo{}
	qpi.QueryIndex = 0
	subqpi := &QueryParseInfo{}
	subqpi.QueryIndex = 1
	qpi.SubEntities = []*QueryParseInfo{subqpi}

	// act
	result := qpi.GetNextQueryIndex()

	// assert
	assert.True(t, result == 2)
}

func TestGetQueryIDRelationMap(t *testing.T) {
	// arrange
	qpi := QueryParseInfo{}
	qpi.QueryIndex = 0

	// act
	result := qpi.GetQueryIDRelationMap(nil)

	// assert
	assert.NotNil(t, result)
}

func TestGetQueryIDRelationMapWithSubEntities(t *testing.T) {
	// arrange
	qpi := QueryParseInfo{}
	qpi.QueryIndex = 0

	subqpi := &QueryParseInfo{}
	subqpi.Parent = &qpi
	subqpi.QueryIndex = 1
	qpi.SubEntities = []*QueryParseInfo{subqpi}

	// act
	result := qpi.GetQueryIDRelationMap(nil)

	// assert
	assert.NotNil(t, result)
	assert.True(t, len(result) == 1)
}

func TestQueryParseInfoParse(t *testing.T) {
	// arrange
	qpi := QueryParseInfo{}
	pqpi := &QueryParseInfo{}
	eo := &godata.ExpandItem{}
	qpi.Init(entities.EntityTypeThing, 0, pqpi, eo)

	// act
	res, err := qpi.Parse(nil)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestCreateWhereIs(t *testing.T) {
	// arrange
	// act
	result := createWhereIs(entities.EntityTypeThing, "field", "prefix")

	// assert
	assert.NotNil(t, result)
	assert.True(t, result == "prefix_")

}

func TestCreateWhereIsWithNoPrefix(t *testing.T) {
	// arrange
	// act
	result := createWhereIs(entities.EntityTypeThing, "id", "")

	// assert
	assert.NotNil(t, result)
	assert.True(t, result == "thing.id")
}

func TestGetJoinForThings(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables
	// act
	resultThingDatastream := getJoin(tables, entities.EntityTypeThing, entities.EntityTypeDatastream, "prefix")
	resultThingHistoricalLocation := getJoin(tables, entities.EntityTypeThing, entities.EntityTypeHistoricalLocation, "prefix")
	resultThingLocation := getJoin(tables, entities.EntityTypeThing, entities.EntityTypeLocation, "prefix")

	// assert
	assert.NotNil(t, resultThingDatastream)
	assert.True(t, resultThingDatastream == "WHERE thing.id = prefix_datastream.datastream_thing_id")
	assert.True(t, resultThingHistoricalLocation == "WHERE thing.id = prefix_historicallocation.historicallocation_thing_id")
	assert.True(t, resultThingLocation == "INNER JOIN  ON thing.id = thing_to_location.thing_id AND prefix_location.location_id = thing_to_location.location_id")
}

func TestGetJoinForLocation(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables
	// act
	resultLocationHistoricalHocation := getJoin(tables, entities.EntityTypeLocation, entities.EntityTypeHistoricalLocation, "prefix")
	resultLocationThing := getJoin(tables, entities.EntityTypeLocation, entities.EntityTypeThing, "prefix")

	// assert
	assert.NotNil(t, resultLocationHistoricalHocation)
	assert.NotNil(t, resultLocationThing)
}

func TestGetJoinForHistoricalLocation(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables
	// act
	resultHistoricalHocationLocation := getJoin(tables, entities.EntityTypeHistoricalLocation, entities.EntityTypeLocation, "prefix")
	resultHistoricalLocationThing := getJoin(tables, entities.EntityTypeHistoricalLocation, entities.EntityTypeThing, "prefix")

	// assert
	assert.NotNil(t, resultHistoricalHocationLocation)
	assert.NotNil(t, resultHistoricalLocationThing)
}

func TestGetJoinForSensor(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	resultSensorDatastream := getJoin(tables, entities.EntityTypeSensor, entities.EntityTypeDatastream, "prefix")

	// assert
	assert.NotNil(t, resultSensorDatastream)
}

func TestGetJoinForObservedProperty(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	resultObservedPropertyDatastream := getJoin(tables, entities.EntityTypeObservedProperty, entities.EntityTypeDatastream, "prefix")

	// assert
	assert.NotNil(t, resultObservedPropertyDatastream)
}

func TestGetJoinForObservation(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	resultObservationDatastream := getJoin(tables, entities.EntityTypeObservation, entities.EntityTypeDatastream, "prefix")
	resultObservationFoi := getJoin(tables, entities.EntityTypeObservation, entities.EntityTypeFeatureOfInterest, "prefix")

	// assert
	assert.NotNil(t, resultObservationDatastream)
	assert.NotNil(t, resultObservationFoi)
}

func TestGetJoinForFoi(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	resultFoiObservation := getJoin(tables, entities.EntityTypeFeatureOfInterest, entities.EntityTypeObservation, "prefix")

	// assert
	assert.NotNil(t, resultFoiObservation)
}

func TestGetJoinForDatastream(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	resultDatastreamThing := getJoin(tables, entities.EntityTypeDatastream, entities.EntityTypeThing, "prefix")
	resultDatastreamSensor := getJoin(tables, entities.EntityTypeDatastream, entities.EntityTypeSensor, "prefix")
	resultDatastreamObservedProperty := getJoin(tables, entities.EntityTypeDatastream, entities.EntityTypeObservedProperty, "prefix")
	resultDatastreamObservation := getJoin(tables, entities.EntityTypeDatastream, entities.EntityTypeObservation, "prefix")
	resultDatastreamLocation := getJoin(tables, entities.EntityTypeDatastream, entities.EntityTypeLocation, "prefix")

	// assert
	assert.NotNil(t, resultDatastreamThing)
	assert.NotNil(t, resultDatastreamSensor)
	assert.NotNil(t, resultDatastreamObservedProperty)
	assert.NotNil(t, resultDatastreamObservation)
	assert.NotNil(t, resultDatastreamLocation)
}

func TestGetJoinForUnknown(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	resultUnknown := getJoin(tables, entities.EntityTypeUnknown, entities.EntityTypeObservation, "prefix")

	// assert
	assert.True(t, resultUnknown == "")
}
