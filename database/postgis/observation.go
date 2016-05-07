package postgis

import (
	"encoding/json"
	"fmt"
	"github.com/geodan/gost/sensorthings/entities"
	"strconv"

	"errors"
	gostErrors "github.com/geodan/gost/errors"
)

// GetObservation todo
func (gdb *GostDatabase) GetObservation(id string) (*entities.Observation, error) {
	////SELECT TO_CHAR(lower(time)AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.Z') AS start, TO_CHAR(upper(time)AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.Z') AS end FROM test;
	return nil, nil
}

// GetObservations todo
func (gdb *GostDatabase) GetObservations() ([]*entities.Observation, error) {
	return nil, nil
}

/*
phenomenonTime	mandatory	String
result	mandatory	String
resultTime	mandatory	String
resultQuality	optional	String
validTime	optional	String
parameters	optional	String


  phenomenontime tstzrange,
  result jsonb,
  resulttime timestamp with time zone,
  resultquality character varying(25),
  validtime tstzrange,
  parameters jsonb,
  stream_id integer,
  featureofinterest_id integer,
*/

// PostObservation todo
func (gdb *GostDatabase) PostObservation(o entities.Observation) (*entities.Observation, error) {
	var oID int

	dID, err := strconv.Atoi(o.Datastream.ID)
	if err != nil || !gdb.DatastreamExists(dID) {
		return nil, gostErrors.NewBadRequestError(errors.New("Datastream does not exist"))
	}

	pTime := PrepareTimeRangeForPostgres(o.PhenomenonTime)
	vTime := PrepareTimeRangeForPostgres(o.ValidTime)
	resultQuality := "NULL"
	parameters := "NULL"
	resultTime := "NULL"
	fID := "NULL"

	if len(o.ResultQuality) != 0 {
		resultQuality = o.ResultQuality
	}

	if len(o.ResultTime) != 0 || o.ResultTime != "NULL" {
		resultTime = fmt.Sprintf("'%s'", o.ResultTime)
	}

	if len(o.Parameters) != 0 {
		pb, _ := json.Marshal(o.Parameters)
		parameters = fmt.Sprintf("%v", pb)
	}

	if o.FeatureOfInterest != nil && len(o.FeatureOfInterest.ID) != 0 {
		fID, err := strconv.Atoi(o.FeatureOfInterest.ID)
		if err != nil || !gdb.FeatureOfInterestExists(fID) {
			return nil, gostErrors.NewBadRequestError(errors.New("FeatureOfInterest does not exist"))
		}
	}

	sql := fmt.Sprintf("INSERT INTO %s.observation (phenomenontime, parameters, validtime, resulttime, result, resultquality, stream_id, featureofinterest_id) VALUES (%s, %s, %s, %s, %v, %s, %v, %v) RETURNING id",
		gdb.Schema, pTime, parameters, vTime, resultTime, o.Result, resultQuality, dID, fID)
	err = gdb.Db.QueryRow(sql).Scan(&oID)
	if err != nil {
		return nil, err
	}

	o.ID = strconv.Itoa(oID)

	// clear inner entities to serves links upon response
	o.Datastream = nil
	o.FeatureOfInterest = nil

	return &o, nil
}
