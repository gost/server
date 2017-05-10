package postgis

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/geodan/gost/src/sensorthings/odata"

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
