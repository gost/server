package postgis

import (
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/odata"
	"github.com/gost/godata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateQueryBuilder(t *testing.T) {
	// act
	qb := CreateQueryBuilder("v1.0", 1)
	// assert
	assert.NotNil(t, qb)
}

func TestRemoveSchema(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)

	// act
	res := qb.removeSchema("v2.hallo")
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

func TestGetOrderByWithQueryOptions(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	qo := &odata.QueryOptions{}
	qob, err := godata.ParseOrderByString("id asc,name desc")
	t.Logf("Error parsing OrderBy string: %v", err)
	qo.OrderBy = qob
	ds := &entities.Datastream{}

	// act
	res := qb.getOrderBy(ds.GetEntityType(), qo)

	// assert
	q := *qob
	t.Logf("count: %v string: %s", len(q.OrderByItems), res)
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

func TestCreateCountQuery(t *testing.T) {
	// arrange
	qb := CreateQueryBuilder("v1.0", 1)
	expected := "SELECT COUNT(*) FROM v1.0.datastream INNER JOIN LATERAL (SELECT thing.id AS thing_id FROM v1.0.thing WHERE thing.id = datastream.thing_id AND thing.id = 1) AS thing on true  WHERE thing.thing_id = 1"
	res := qb.CreateCountQuery(&entities.Datastream{}, &entities.Thing{}, 1, nil)

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
