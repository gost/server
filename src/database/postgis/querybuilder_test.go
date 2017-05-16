package postgis

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/geodan/gost/src/sensorthings/odata"
	"github.com/geodan/gost/src/sensorthings/entities"


)

func TestCreateQueryBuilder(t *testing.T){
	// act
	qb := CreateQueryBuilder("v1.0",1)
	// assert
	assert.NotNil(t, qb)
}

func TestRemoveSchema(t *testing.T){
	// arrange
	qb := CreateQueryBuilder("v1.0",1)

	// act
	res := qb.removeSchema("v2.hallo")
	// assert
	assert.True(t, res=="hallo")
}

func TestGetLimit(t *testing.T){
	// arrange
	qb := CreateQueryBuilder("v1.0",1)
	qo := &odata.QueryOptions{}

	// act
	res := qb.getLimit(qo)
	// assert
	assert.True(t, res=="1")
}

func TestGetLimitWithQueryTop(t *testing.T){
	// arrange
	qb := CreateQueryBuilder("v1.0",1)
	qo := &odata.QueryOptions{}
	qo.QueryTop = &odata.QueryTop{odata.QueryBase{"0"},2}

	// act
	res := qb.getLimit(qo)
	// assert
	assert.True(t, res=="2")
}

func TestCreateCountQuery(t *testing.T){
	// arrange
	qb := CreateQueryBuilder("v1.0",1)
	expected := "SELECT COUNT(*) FROM v1.0.datastream INNER JOIN LATERAL (SELECT thing.id AS thing_id FROM v1.0.thing WHERE thing.id = datastream.thing_id ) AS thing on true  WHERE thing.thing_id = 1"

	// act
	res := qb.CreateCountQuery(&entities.Datastream{}, &entities.Thing{}, 1, nil)

	// assert
	assert.NotNil(t, res)
	assert.True(t, expected == res)
}

func TestCreateQuery(t *testing.T){
	// arrange
	qb := CreateQueryBuilder("v1.0",1)
	expected := "SELECT datastream.id AS A_datastream_id, datastream.name AS A_datastream_name, datastream.description AS A_datastream_description, datastream.unitofmeasurement AS A_datastream_unitofmeasurement, datastream.observationtype AS A_datastream_observationtype, public.ST_AsGeoJSON(datastream.observedarea) AS A_datastream_observedarea, datastream.phenomenontime AS A_datastream_phenomenontime, datastream.resulttime AS A_datastream_resulttime FROM v1.0.datastream  WHERE datastream.id IN (SELECT datastream.id FROM v1.0.datastream  WHERE datastream.id = 0 ORDER BY datastream.id DESC  OFFSET 0) ORDER BY datastream.id DESC"

	// act
	query, _ := qb.CreateQuery(&entities.Datastream{}, nil, 0, nil)

	// assert
	assert.NotNil(t, query)
	assert.True(t, expected == query)
}