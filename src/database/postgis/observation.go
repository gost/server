package postgis

import (
	"errors"
	"fmt"
	"strings"

	"database/sql"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
)

var totalObservations int

// GetTotalObservations returns the amount of observations in the database
func (gdb *GostDatabase) GetTotalObservations() int {
	return totalObservations
}

// InitObservations Initialises the datastream repository, setting totalObservations on startup
func (gdb *GostDatabase) InitObservations() {
	sql := fmt.Sprintf("SELECT Count(*) from %s.observation", gdb.Schema)
	gdb.Db.QueryRow(sql).Scan(&totalObservations)
}

// GetObservation retrieves an observation by id from the database
func (gdb *GostDatabase) GetObservation(id interface{}, qo *odata.QueryOptions) (*entities.Observation, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Observation does not exist"))
	}

	sql := fmt.Sprintf("select id, data FROM %s.observation where id = %v ", gdb.Schema, intID)
	observation, err := processObservation(gdb.Db, sql, qo)
	if err != nil {
		return nil, err
	}

	return observation, nil
}

// GetObservations retrieves all datastreams
func (gdb *GostDatabase) GetObservations(qo *odata.QueryOptions) ([]*entities.Observation, error) {
	sql := fmt.Sprintf("select id, data FROM %s.observation order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema)
	return processObservations(gdb.Db, sql, qo)
}

// GetObservationsByFeatureOfInterest retrieves all observations by the given FeatureOfInterest id
func (gdb *GostDatabase) GetObservationsByFeatureOfInterest(foiID interface{}, qo *odata.QueryOptions) ([]*entities.Observation, error) {
	intID, ok := ToIntID(foiID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest does not exist"))
	}

	sql := fmt.Sprintf("select id, data FROM %s.observation where featureofinterest_id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, intID)
	return processObservations(gdb.Db, sql, qo)
}

// GetObservationsByDatastream retrieves all observations by the given datastream id
func (gdb *GostDatabase) GetObservationsByDatastream(dataStreamID interface{}, qo *odata.QueryOptions) ([]*entities.Observation, error) {
	intID, ok := ToIntID(dataStreamID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select id, data FROM %s.observation where stream_id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, intID)
	return processObservations(gdb.Db, sql, qo)
}

func processObservation(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.Observation, error) {
	observations, err := processObservations(db, sql, qo)
	if err != nil {
		return nil, err
	}

	if len(observations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Observation not found"))
	}

	return observations[0], nil
}

func processObservations(db *sql.DB, sql string, qo *odata.QueryOptions) ([]*entities.Observation, error) {
	rows, err := db.Query(sql)
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
		observation.ID = id
		err = observation.ParseEntity([]byte(data))

		if err != nil {
			return nil, err
		}

		if qo != nil && qo.QuerySelect != nil && len(qo.QuerySelect.Params) > 0 {
			set := make(map[string]bool)
			for _, v := range qo.QuerySelect.Params {
				set[v] = true
			}

			_, ok := set["id"]
			if !ok {
				observation.ID = nil
			}
			_, ok = set["phenomenonTime"]
			if !ok {
				observation.PhenomenonTime = ""
			}
			_, ok = set["result"]
			if !ok {
				observation.Result = nil
			}
			_, ok = set["resultTime"]
			if !ok {
				observation.ResultTime = ""
			}
			_, ok = set["resultQuality"]
			if !ok {
				observation.ResultQuality = ""
			}
			_, ok = set["validTime"]
			if !ok {
				observation.ValidTime = ""
			}
			_, ok = set["parameters"]
			if !ok {
				observation.Parameters = nil
			}
		}

		observations = append(observations, &observation)
	}

	return observations, nil
}

// PostObservation todo
func (gdb *GostDatabase) PostObservation(o *entities.Observation) (*entities.Observation, error) {
	var oID int

	dID, ok := ToIntID(o.Datastream.ID)
	if !ok {
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

	err := gdb.Db.QueryRow(sql).Scan(&oID)
	if err != nil {
		errString := fmt.Sprintf("%v", err.Error())
		if strings.Contains(errString, "violates foreign key constraint \"fk_datastream\"") {
			return nil, gostErrors.NewBadRequestError(err)
		}

		return nil, err
	}

	o.ID = oID
	if o.ResultTime == "NULL" {
		o.ResultTime = ""
	}

	// clear inner entities to serves links upon response
	o.Datastream = nil
	o.FeatureOfInterest = nil

	totalObservations++
	return o, nil
}

// DeleteObservation tries to delete a Observation by the given id
func (gdb *GostDatabase) DeleteObservation(id interface{}) error {
	intID, ok := ToIntID(id)
	if !ok {
		return gostErrors.NewRequestNotFound(errors.New("Observation does not exist"))
	}

	r, err := gdb.Db.Exec(fmt.Sprintf("DELETE FROM %s.observation WHERE id = $1", gdb.Schema), intID)
	if err != nil {
		return err
	}

	if c, _ := r.RowsAffected(); c == 0 {
		return gostErrors.NewRequestNotFound(errors.New("Observation not found"))
	}

	totalObservations--
	return nil
}
