package postgis

import (
	"encoding/json"
	"fmt"
	"github.com/geodan/gost/sensorthings/entities"
	"strconv"

	"database/sql"
	"errors"
	gostErrors "github.com/geodan/gost/errors"
)

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

// GetObservation todo
func (gdb *GostDatabase) GetObservation(id string) (*entities.Observation, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select id, to_char(lower(phenomenontime) at time zone 'UTC', '%s') as phenomenontimeStart, "+
		"to_char(upper(phenomenontime) at time zone 'UTC', '%s') as phenomenontimeEnd, "+
		"result, to_char(resulttime at time zone 'UTC', '%s') as resulttime, to_char(lower(validtime) at time zone 'UTC', '%s') as validtimeStart, "+
		"to_char(upper(validtime) at time zone 'UTC', '%s') as validtimeEnd, resultquality, parameters FROM observation where id = $1", TimeFormat, TimeFormat, TimeFormat, TimeFormat, TimeFormat)

	datastream, err := processObservation(gdb.Db, sql, intID)
	if err != nil {
		return nil, err
	}

	return datastream, nil
}

// GetObservations retrieves all datastreams
func (gdb *GostDatabase) GetObservations() ([]*entities.Observation, error) {
	sql := fmt.Sprintf("select id, to_char(lower(phenomenontime) at time zone 'UTC', '%s') as phenomenontimeStart, "+
		"to_char(upper(phenomenontime) at time zone 'UTC', '%s') as phenomenontimeEnd, "+
		"result, to_char(resulttime at time zone 'UTC', '%s') as resulttime, to_char(lower(validtime) at time zone 'UTC', '%s') as validTimeStart, "+
		"to_char(upper(validtime) at time zone 'UTC', '%s') as validTimeEnd, resultquality, parameters FROM observation", TimeFormat, TimeFormat, TimeFormat, TimeFormat, TimeFormat)
	return processObservations(gdb.Db, sql)
}

func processObservation(db *sql.DB, sql string, args ...interface{}) (*entities.Observation, error) {
	observations, err := processObservations(db, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(observations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Observation not found"))
	}

	return observations[0], nil
}

func processObservations(db *sql.DB, sql string, args ...interface{}) ([]*entities.Observation, error) {
	rows, err := db.Query(sql, args...)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var observations = []*entities.Observation{}
	for rows.Next() {
		var id int
		var result float64
		var phenomenonTimeStart, phenomenonTimeEnd, validTimeStart, validTimeEnd, parameters, resultQuality, resultTime *string

		err := rows.Scan(&id, &phenomenonTimeStart, &phenomenonTimeEnd, &result, &resultTime, &validTimeStart, &validTimeEnd, &resultQuality, &parameters)
		if err != nil {
			return nil, err
		}

		parametersMap, err := JSONToMap(parameters)

		if err != nil {
			return nil, err
		}

		datastream := entities.Observation{
			ID:             strconv.Itoa(id),
			ResultTime:     ConvertNullString(resultTime),
			Result:         result,
			PhenomenonTime: TimeRangeToString(phenomenonTimeStart, phenomenonTimeEnd),
			ValidTime:      TimeRangeToString(validTimeStart, validTimeEnd),
			Parameters:     parametersMap,
			ResultQuality:  ConvertNullString(resultQuality),
		}
		observations = append(observations, &datastream)

	}

	return observations, nil
}

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
		resultQuality = fmt.Sprintf("'%s'", o.ResultQuality)
	}

	if len(o.ResultTime) != 0 && o.ResultTime != "NULL" {
		resultTime = fmt.Sprintf("'%s'", o.ResultTime)
	}

	if len(o.Parameters) != 0 {
		pb, _ := json.Marshal(o.Parameters)
		parameters = fmt.Sprintf("'%s'", string(pb[:]))
	}

	if o.FeatureOfInterest != nil && len(o.FeatureOfInterest.ID) != 0 {
		fID, err := strconv.Atoi(o.FeatureOfInterest.ID)
		if err != nil || !gdb.FeatureOfInterestExists(fID) {
			return nil, gostErrors.NewBadRequestError(errors.New("FeatureOfInterest does not exist"))
		}
	}

	sql := fmt.Sprintf("INSERT INTO observation (phenomenontime, parameters, validtime, resulttime, result, resultquality, stream_id, featureofinterest_id) VALUES (%s, %v, %s, %s, %v, %s, %v, %v) RETURNING id",
		pTime, parameters, vTime, resultTime, o.Result, resultQuality, dID, fID)

	fmt.Println(sql)
	err = gdb.Db.QueryRow(sql).Scan(&oID)
	if err != nil {
		return nil, err
	}

	o.ID = strconv.Itoa(oID)

	if o.ResultTime == "NULL" {
		o.ResultTime = ""
	}

	// clear inner entities to serves links upon response
	o.Datastream = nil
	o.FeatureOfInterest = nil

	return &o, nil
}
