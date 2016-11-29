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

// observationParamFactory is used to construct a WHERE clause from an ODATA $select string
func observationParamFactory(key string, value interface{}) (string, string, error) {
	val := fmt.Sprintf("%v", value)
	jsonVal := convertSelectValueForJSON(val)
	switch key {
	case "id":
		return "id", val, nil
		break
	case "phenomenonTime":
		return "data -> 'phenomenonTime'", jsonVal, nil
		break
	case "resultTime":
		return "data -> 'resultTime'", jsonVal, nil
		break
	case "result":
		return "data -> 'result'", jsonVal, nil
		break
	case "resultQuality":
		return "data -> 'resultQuality'", jsonVal, nil
		break
	case "parameters": //implement parameters/parameterName
		return "data -> 'parameters'", jsonVal, nil
		break
	}

	return "", "", fmt.Errorf("Parameter %s not implemented", key)
}

func convertSelectValueForJSON(value string) string {
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		return value
	}

	return fmt.Sprintf("'%v'", value)
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
func (gdb *GostDatabase) GetObservations(qo *odata.QueryOptions) ([]*entities.Observation, int, error) {
	var queryString string
	var err error
	if queryString, err = CreateFilterQueryString(qo, observationParamFactory, "WHERE "); err != nil {
		return nil, 0, gostErrors.NewBadRequestError(err)
	}

	sql := fmt.Sprintf("select id, data FROM %s.observation "+queryString+"order by id desc"+CreateTopSkipQueryString(qo), gdb.Schema)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.observation", gdb.Schema)
	return processObservations(gdb.Db, sql, qo, countSQL)
}

// GetObservationsByFeatureOfInterest retrieves all observations by the given FeatureOfInterest id
func (gdb *GostDatabase) GetObservationsByFeatureOfInterest(foiID interface{}, qo *odata.QueryOptions) ([]*entities.Observation, int, error) {
	intID, ok := ToIntID(foiID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest does not exist"))
	}

	sql := fmt.Sprintf("select id, data FROM %s.observation where featureofinterest_id = %v order by id desc"+CreateTopSkipQueryString(qo), gdb.Schema, intID)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.observation where featureofinterest_id = %v", gdb.Schema, intID)
	return processObservations(gdb.Db, sql, qo, countSQL)
}

// GetObservationsByDatastream retrieves all observations by the given datastream id
func (gdb *GostDatabase) GetObservationsByDatastream(dataStreamID interface{}, qo *odata.QueryOptions) ([]*entities.Observation, int, error) {
	var intID int
	var queryString string
	var ok bool
	var err error

	if intID, ok = ToIntID(dataStreamID); !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	if queryString, err = CreateFilterQueryString(qo, observationParamFactory, " AND "); err != nil {
		return nil, 0, gostErrors.NewBadRequestError(errors.New("Datastream does not exist"))
	}

	sql := fmt.Sprintf("select id, data FROM %s.observation where stream_id = %v "+queryString+"order by id desc"+CreateTopSkipQueryString(qo), gdb.Schema, intID)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.observation where stream_id = %v", gdb.Schema, intID)
	return processObservations(gdb.Db, sql, qo, countSQL)
}

func processObservation(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.Observation, error) {
	observations, _, err := processObservations(db, sql, qo, "")
	if err != nil {
		return nil, err
	}

	if len(observations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Observation not found"))
	}

	return observations[0], nil
}

func processObservations(db *sql.DB, sql string, qo *odata.QueryOptions, countSQL string) ([]*entities.Observation, int, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	var observations = []*entities.Observation{}
	for rows.Next() {
		var id int
		var data string

		err := rows.Scan(&id, &data)
		if err != nil {
			return nil, 0, err
		}

		observation := entities.Observation{}
		observation.ID = id
		err = observation.ParseEntity([]byte(data))

		if err != nil {
			return nil, 0, err
		}

		if qo != nil && qo.QuerySelect != nil && len(qo.QuerySelect.Params) > 0 {
			set := make(map[string]bool)
			for _, v := range qo.QuerySelect.Params {
				set[strings.ToLower(v)] = true
			}

			_, ok := set["id"]
			if !ok {
				observation.ID = nil
			}
			_, ok = set["phenomenontime"]
			if !ok {
				observation.PhenomenonTime = ""
			}
			_, ok = set["result"]
			if !ok {
				observation.Result = nil
			}
			_, ok = set["resulttime"]
			if !ok {
				observation.ResultTime = ""
			}
			_, ok = set["resultquality"]
			if !ok {
				observation.ResultQuality = ""
			}
			_, ok = set["validtime"]
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

	var count int
	if len(countSQL) > 0 {
		db.QueryRow(countSQL).Scan(&count)
	}

	return observations, count, nil
}

// PutObservation replaces an observation to the database
func (gdb *GostDatabase) PutObservation(id interface{}, o *entities.Observation) (*entities.Observation, error) {
	return gdb.PatchObservation(id, o)
	/*var err error
	var ok bool
	var intID, dID, fID int

	if intID, ok = ToIntID(id); !ok || !gdb.ObservationExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Observation does not exist"))
	}

	if dID, ok = ToIntID(o.Datastream.ID); !ok || !gdb.ObservationExists(dID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	if fID, ok = ToIntID(o.FeatureOfInterest.ID); !ok || !gdb.FeatureOfInterestExists(fID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest does not exist"))
	}

	json, _ := o.MarshalPostgresJSON()
	obs := fmt.Sprintf("'%s'", string(json[:]))
	sql := fmt.Sprintf("update %s.observation set data=%v, stream_id=%v, featureofinterest_id=%v where id=%v", gdb.Schema, obs, dID, fID, intID)

	_, err2 := gdb.Db.Exec(sql)
	if err2 != nil {
		return nil, err
	}

	o.ID = intID
	return o, nil
	*/
}

// PostObservation adds an observation to the database
func (gdb *GostDatabase) PostObservation(o *entities.Observation) (*entities.Observation, error) {
	var oID int

	dID, ok := ToIntID(o.Datastream.ID)
	if !ok {
		return nil, gostErrors.NewBadRequestError(errors.New("Datastream does not exist"))
	}

	if o.FeatureOfInterest == nil || len(fmt.Sprintf("%v", o.FeatureOfInterest.ID)) == 0 {
		return nil, gostErrors.NewBadRequestError(errors.New("No FeatureOfInterest supplied or Location found on linked thing"))
	}

	fID := o.FeatureOfInterest.ID

	json, _ := o.MarshalPostgresJSON()
	obs := fmt.Sprintf("'%s'", string(json[:]))
	sql := fmt.Sprintf("INSERT INTO %s.observation (data, stream_id, featureofinterest_id) VALUES (%v, %v, %v) RETURNING id", gdb.Schema, obs, dID, fID)

	err := gdb.Db.QueryRow(sql).Scan(&oID)
	if err != nil {
		errString := fmt.Sprintf("%v", err.Error())
		if strings.Contains(errString, "violates foreign key constraint \"fk_datastream\"") {
			return nil, gostErrors.NewBadRequestError(errors.New("Datastream does not exist"))
		}
		if strings.Contains(errString, "violates foreign key constraint \"fk_featureofinterest\"") {
			return nil, gostErrors.NewBadRequestError(errors.New("FeatureOfInterest does not exist"))
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

	return o, nil
}

// ObservationExists checks if an Observation is present in the database based on a given id.
func (gdb *GostDatabase) ObservationExists(id interface{}) bool {
	return EntityExists(gdb, id, "observation")
}

// PatchObservation updates a Observation in the database
func (gdb *GostDatabase) PatchObservation(id interface{}, o *entities.Observation) (*entities.Observation, error) {
	var err error
	var ok bool
	var intID int
	updates := make(map[string]interface{})

	if intID, ok = ToIntID(id); !ok || !gdb.ObservationExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Observation does not exist"))
	}

	observation, _ := gdb.GetObservation(intID, nil)

	if len(o.PhenomenonTime) > 0 {
		observation.PhenomenonTime = o.PhenomenonTime
	}

	if o.Result != nil {
		observation.Result = o.Result
	}

	if len(o.ResultTime) > 0 {
		observation.ResultTime = o.ResultTime
	}

	if len(o.ResultQuality) > 0 {
		observation.ResultQuality = o.ResultQuality
	}

	if len(o.ValidTime) > 0 {
		observation.ValidTime = o.ValidTime
	}

	if len(o.Parameters) > 0 {
		observation.Parameters = o.Parameters
	}

	json, _ := observation.MarshalPostgresJSON()
	updates["data"] = string(json[:])

	if err = gdb.updateEntityColumns("observation", updates, intID); err != nil {
		return nil, err
	}

	return observation, nil
}

// DeleteObservation tries to delete a Observation by the given id
func (gdb *GostDatabase) DeleteObservation(id interface{}) error {
	return DeleteEntity(gdb, id, "observation")
}
