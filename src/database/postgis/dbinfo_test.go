package postgis

import (
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
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

	eo := &odata.ExpandOperation{}

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
	qpi_return := qpi.GetQueryParseInfoByQueryIndex(0)

	// assert
	assert.NotNil(t, qpi_return)
	assert.NotNil(t, qpi_return.QueryIndex == 0)

}

func TestGetQueryParseInfoByQueryIndexWithSubEntities(t *testing.T) {
	// arrange
	qpi := QueryParseInfo{}
	subqpi := &QueryParseInfo{}
	subqpi.QueryIndex = 1
	qpi.SubEntities = []*QueryParseInfo{subqpi}

	// act
	qpi_return := qpi.GetQueryParseInfoByQueryIndex(1)

	// assert
	assert.True(t, qpi_return.QueryIndex == 1)
}

func TestGetQueryParseInfoByQueryIndexNotFound(t *testing.T) {
	// arrange
	qpi := QueryParseInfo{}
	subqpi := &QueryParseInfo{}
	subqpi.QueryIndex = 1
	qpi.SubEntities = []*QueryParseInfo{subqpi}

	// act
	qpi_return := qpi.GetQueryParseInfoByQueryIndex(2)

	// assert
	assert.Nil(t, qpi_return)
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
	eo := &odata.ExpandOperation{}
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
	result_thing_datastream := getJoin(tables, entities.EntityTypeThing, entities.EntityTypeDatastream, "prefix")
	result_thing_historicallocation := getJoin(tables, entities.EntityTypeThing, entities.EntityTypeHistoricalLocation, "prefix")
	result_thing_location := getJoin(tables, entities.EntityTypeThing, entities.EntityTypeLocation, "prefix")

	// assert
	assert.NotNil(t, result_thing_datastream)
	assert.True(t, result_thing_datastream == "WHERE thing.id = prefix_datastream.datastream_thing_id")
	assert.True(t, result_thing_historicallocation == "WHERE thing.id = prefix_historicallocation.historicallocation_thing_id")
	assert.True(t, result_thing_location == "INNER JOIN  ON thing.id = thing_to_location.thing_id AND prefix_location.location_id = thing_to_location.location_id")
}

func TestGetJoinForLocation(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables
	// act
	result_location_historicallocation := getJoin(tables, entities.EntityTypeLocation, entities.EntityTypeHistoricalLocation, "prefix")
	result_location_thing := getJoin(tables, entities.EntityTypeLocation, entities.EntityTypeThing, "prefix")

	// assert
	assert.NotNil(t, result_location_historicallocation)
	assert.NotNil(t, result_location_thing)
}

func TestGetJoinForHistoricalLocation(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables
	// act
	result_historicallocation_location := getJoin(tables, entities.EntityTypeHistoricalLocation, entities.EntityTypeLocation, "prefix")
	result_historicallocation_thing := getJoin(tables, entities.EntityTypeHistoricalLocation, entities.EntityTypeThing, "prefix")

	// assert
	assert.NotNil(t, result_historicallocation_location)
	assert.NotNil(t, result_historicallocation_thing)
}

func TestGetJoinForSensor(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	result_sensor_datastream := getJoin(tables, entities.EntityTypeSensor, entities.EntityTypeDatastream, "prefix")

	// assert
	assert.NotNil(t, result_sensor_datastream)
}

func TestGetJoinForObservedProperty(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	result_observedproperty_datastream := getJoin(tables, entities.EntityTypeObservedProperty, entities.EntityTypeDatastream, "prefix")

	// assert
	assert.NotNil(t, result_observedproperty_datastream)
}

func TestGetJoinForObservation(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	result_observation_datastream := getJoin(tables, entities.EntityTypeObservation, entities.EntityTypeDatastream, "prefix")
	result_observation_foi := getJoin(tables, entities.EntityTypeObservation, entities.EntityTypeFeatureOfInterest, "prefix")

	// assert
	assert.NotNil(t, result_observation_datastream)
	assert.NotNil(t, result_observation_foi)
}

func TestGetJoinForFoi(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	result_foi_observation := getJoin(tables, entities.EntityTypeFeatureOfInterest, entities.EntityTypeObservation, "prefix")

	// assert
	assert.NotNil(t, result_foi_observation)
}


func TestGetJoinForDatastream(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	result_datastream_thing := getJoin(tables, entities.EntityTypeDatastream, entities.EntityTypeThing, "prefix")
	result_datastream_sensor := getJoin(tables, entities.EntityTypeDatastream, entities.EntityTypeSensor, "prefix")
	result_datastream_observedproperty := getJoin(tables, entities.EntityTypeDatastream, entities.EntityTypeObservedProperty, "prefix")
	result_datastream_observation := getJoin(tables, entities.EntityTypeDatastream, entities.EntityTypeObservation, "prefix")
	result_datastream_location := getJoin(tables, entities.EntityTypeDatastream, entities.EntityTypeLocation, "prefix")

	// assert
	assert.NotNil(t, result_datastream_thing)
	assert.NotNil(t, result_datastream_sensor)
	assert.NotNil(t, result_datastream_observedproperty)
	assert.NotNil(t, result_datastream_observation)
	assert.NotNil(t, result_datastream_location)
}


func TestGetJoinForUnknown(t *testing.T) {
	// arrange
	qb := QueryBuilder{}
	tables := qb.tables

	// act
	result_unknown := getJoin(tables, entities.EntityTypeUnknown, entities.EntityTypeObservation, "prefix")

	// assert
	assert.True(t, result_unknown=="")
}
