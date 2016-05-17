package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"strconv"

	"database/sql"
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
)

// GetObservation todo
func (gdb *GostDatabase) GetObservation(id string) (*entities.Observation, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select id, data FROM %s.observation where id = $1", gdb.Schema)

	datastream, err := processObservation(gdb.Db, sql, intID)
	if err != nil {
		return nil, err
	}

	return datastream, nil
}

// GetObservations retrieves all datastreams
func (gdb *GostDatabase) GetObservations() ([]*entities.Observation, error) {
	sql := fmt.Sprintf("select id, data FROM %s.observation", gdb.Schema)
	return processObservations(gdb.Db, sql)
}

// GetObservationsByDatastream retrieves all observations by the given datastream id
func (gdb *GostDatabase) GetObservationsByDatastream(dataStreamID string) ([]*entities.Observation, error) {
	intID, err := strconv.Atoi(dataStreamID)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select id, data FROM %s.observation where stream_id = $1", gdb.Schema)
	return processObservations(gdb.Db, sql, intID)
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
		var data string

		err := rows.Scan(&id, &data)
		if err != nil {
			return nil, err
		}

		observation := entities.Observation{}
		observation.ID = strconv.Itoa(id)
		err = observation.ParseEntity([]byte(data))

		if err != nil {
			return nil, err
		}

		observations = append(observations, &observation)
	}

	return observations, nil
}

// PostObservation todo
func (gdb *GostDatabase) PostObservation(o entities.Observation) (*entities.Observation, error) {
	var oID int

	dID, err := strconv.Atoi(o.Datastream.ID)
	if err != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("Datastream does not exist"))
	}

	fID := "NULL"

	/*if o.FeatureOfInterest != nil && len(o.FeatureOfInterest.ID) != 0 {
		fID, err := strconv.Atoi(o.FeatureOfInterest.ID)
		if err != nil {
			return nil, gostErrors.NewBadRequestError(errors.New("FeatureOfInterest does not exist"))
		}
	}*/
	//REMOVE EXIST and insert return error based on database error

	//marshal for
	json, _ := o.MarshalPostgresJSON()
	obs := fmt.Sprintf("'%s'", string(json[:]))
	sql := fmt.Sprintf("INSERT INTO %s.observation (data, stream_id, featureofinterest_id) VALUES (%v, %v, %v) RETURNING id", gdb.Schema, obs, dID, fID)

	//ToDo: Check error fk exist?
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
