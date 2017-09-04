package postgis

import (
	"errors"
	"fmt"
	"strings"

	"database/sql"

	gostErrors "github.com/gost/server/errors"
	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/odata"
	"strconv"
)

func observationParamFactory(values map[string]interface{}) (entities.Entity, error) {
	o := &entities.Observation{}
	for as, value := range values {
		if as == asMappings[entities.EntityTypeObservation][observationResultTime] {
			if value == nil {
				empty := ""
				o.ResultTime = &empty
			} else {
				rt := value.(string)
				o.ResultTime = &rt
			}
		}

		if value == nil {
			continue
		}
		if as == asMappings[entities.EntityTypeObservation][observationID] {
			o.ID = value
		}
		if as == asMappings[entities.EntityTypeObservation][observationPhenomenonTime] {
			o.PhenomenonTime = value.(string)
		}
		if as == asMappings[entities.EntityTypeObservation][observationResult] {
			var result interface{}
			var err error
			if strings.HasPrefix(value.(string), "\"") {
				result = strings.Replace(value.(string), "\"", "", 2)
			} else {
				result, err = strconv.ParseFloat(value.(string), 64)
				if err != nil {
					result = strings.Replace(value.(string), "\"", "", 2)
				}
			}

			o.Result = result
		}
		if as == asMappings[entities.EntityTypeObservation][observationValidTime] {
			o.ValidTime = value.(string)
		}
		if as == asMappings[entities.EntityTypeObservation][observationResultQuality] {
			o.ResultQuality = value.(string)
		}
		if as == asMappings[entities.EntityTypeObservation][observationParameters] {
			t := value.(string)
			parameterMap, err := JSONToMap(&t)
			if err != nil {
				return nil, err
			}

			o.Parameters = parameterMap
		}
	}

	return o, nil
}

// GetObservation retrieves an observation by id from the database
func (gdb *GostDatabase) GetObservation(id interface{}, qo *odata.QueryOptions) (*entities.Observation, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Observation does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Observation{}, nil, intID, qo)
	observation, err := processObservation(gdb.Db, query, qi)
	if err != nil {
		return nil, err
	}

	return observation, nil
}

// GetObservations retrieves all observations
func (gdb *GostDatabase) GetObservations(qo *odata.QueryOptions) ([]*entities.Observation, int, error) {
	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Observation{}, nil, nil, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Observation{}, nil, nil, qo)
	return processObservations(gdb.Db, query, qi, countSQL)
}

// GetObservationsByFeatureOfInterest retrieves all observations by the given FeatureOfInterest id
func (gdb *GostDatabase) GetObservationsByFeatureOfInterest(foiID interface{}, qo *odata.QueryOptions) ([]*entities.Observation, int, error) {
	intID, ok := ToIntID(foiID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Observation{}, &entities.FeatureOfInterest{}, intID, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Observation{}, &entities.FeatureOfInterest{}, intID, qo)
	return processObservations(gdb.Db, query, qi, countSQL)
}

// GetObservationsByDatastream retrieves all observations by the given datastream id
func (gdb *GostDatabase) GetObservationsByDatastream(dataStreamID interface{}, qo *odata.QueryOptions) ([]*entities.Observation, int, error) {
	intID, ok := ToIntID(dataStreamID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.Observation{}, &entities.Datastream{}, intID, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.Observation{}, &entities.Datastream{}, intID, qo)
	return processObservations(gdb.Db, query, qi, countSQL)
}

func processObservation(db *sql.DB, sql string, qi *QueryParseInfo) (*entities.Observation, error) {
	observations, _, err := processObservations(db, sql, qi, "")
	if err != nil {
		return nil, err
	}

	if len(observations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("Observation not found"))
	}

	return observations[0], nil
}

func processObservations(db *sql.DB, sql string, qi *QueryParseInfo, countSQL string) ([]*entities.Observation, int, error) {
	data, err := ExecuteSelect(db, qi, sql)
	if err != nil {
		return nil, 0, fmt.Errorf("Error executing query %v", err)
	}

	o := make([]*entities.Observation, 0)
	for _, d := range data {
		entity := d.(*entities.Observation)
		o = append(o, entity)
	}

	var count int
	if len(countSQL) > 0 {
		count, err = ExecuteSelectCount(db, countSQL)
		if err != nil {
			return nil, 0, fmt.Errorf("Error executing count %v", err)
		}
	}

	return o, count, nil
}

// PutObservation replaces an observation to the database
func (gdb *GostDatabase) PutObservation(id interface{}, o *entities.Observation) (*entities.Observation, error) {
	return gdb.PatchObservation(id, o)
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
	sql2 := fmt.Sprintf("INSERT INTO %s.observation (data, stream_id, featureofinterest_id) VALUES (%v, %v, %v) RETURNING id", gdb.Schema, obs, dID, fID)

	err := gdb.Db.QueryRow(sql2).Scan(&oID)
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

	if o.ResultTime != nil {
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
