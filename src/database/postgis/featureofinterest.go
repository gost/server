package postgis

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
	"strconv"
)

var foiMapping = map[string]string{"feature": "public.ST_AsGeoJSON(featureofinterest.feature) AS feature"}

// GetFeatureOfInterest returns a feature of interest by id
func (gdb *GostDatabase) GetFeatureOfInterest(id interface{}, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.FeatureOfInterest{}, qo, "", "", foiMapping)+" from %s.featureofinterest where id = %v", gdb.Schema, intID)
	return processFeatureOfInterest(gdb.Db, sql, qo)
}

// GetFeatureOfInterestByObservation returns a feature of interest by given observation id
func (gdb *GostDatabase) GetFeatureOfInterestByObservation(id interface{}, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.FeatureOfInterest{}, qo, "", "", foiMapping)+" from %s.featureofinterest inner join %s.observation on observation.featureofinterest_id = featureofinterest.id where observation.id = %v limit 1", gdb.Schema, gdb.Schema, intID)
	return processFeatureOfInterest(gdb.Db, sql, qo)
}

// GetFeatureOfInterests returns all feature of interests
func (gdb *GostDatabase) GetFeatureOfInterests(qo *odata.QueryOptions) ([]*entities.FeatureOfInterest, error) {
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.FeatureOfInterest{}, qo, "", "", foiMapping)+" from %s.featureofinterest", gdb.Schema)
	return processFeatureOfInterests(gdb.Db, sql, qo)
}

// PostFeatureOfInterest inserts a new FeatureOfInterest into the database
func (gdb *GostDatabase) PostFeatureOfInterest(f *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	var fID int
	locationBytes, _ := json.Marshal(f.Feature)
	encoding, _ := entities.CreateEncodingType(f.EncodingType)
	sql := fmt.Sprintf("INSERT INTO %s.featureofinterest (description, encodingtype, feature) VALUES ($1, $2, public.ST_GeomFromGeoJSON('%s')) RETURNING id", gdb.Schema, string(locationBytes[:]))
	err := gdb.Db.QueryRow(sql, f.Description, encoding.Code).Scan(&fID)
	if err != nil {
		return nil, err
	}

	f.ID = strconv.Itoa(fID)
	return f, nil
}

func processFeatureOfInterest(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error) {
	locations, err := processFeatureOfInterests(db, sql, qo)
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest not found"))
	}

	return locations[0], nil
}

func processFeatureOfInterests(db *sql.DB, sql string, qo *odata.QueryOptions) ([]*entities.FeatureOfInterest, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var featureOfInterests = []*entities.FeatureOfInterest{}
	for rows.Next() {
		var ID interface{}
		var encodingType int
		var description, feature string

		var params []interface{}
		var qp []string
		if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
			f := &entities.FeatureOfInterest{}
			qp = f.GetPropertyNames()
		} else {
			qp = qo.QuerySelect.Params
		}

		for _, p := range qp {
			if p == "id" {
				params = append(params, &ID)
			}
			if p == "encodingtype" {
				params = append(params, &encodingType)
			}
			if p == "description" {
				params = append(params, &description)
			}
			if p == "feature" {
				params = append(params, &feature)
			}
		}

		err = rows.Scan(params...)

		featureMap, err := JSONToMap(&feature)
		if err != nil {
			return nil, err
		}

		foi := entities.FeatureOfInterest{}
		foi.ID = ID
		foi.Description = description
		foi.Feature = featureMap
		if encodingType != 0 {
			foi.EncodingType = entities.EncodingValues[encodingType].Value
		}

		featureOfInterests = append(featureOfInterests, &foi)
	}

	return featureOfInterests, nil
}

// DeleteFeatureOfInterest tries to delete a FeatureOfInterest by the given id
func (gdb *GostDatabase) DeleteFeatureOfInterest(id interface{}) error {
	intID, ok := ToIntID(id)
	if !ok {
		return gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest does not exist"))
	}

	r, err := gdb.Db.Exec(fmt.Sprintf("DELETE FROM %s.featureofinterest WHERE id = $1", gdb.Schema), intID)
	if err != nil {
		return err
	}

	if c, _ := r.RowsAffected(); c == 0 {
		return gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest not found"))
	}

	return nil
}

// FeatureOfInterestExists checks if a FeatureOfInterest is present in the database based on a given id.
func (gdb *GostDatabase) FeatureOfInterestExists(databaseID int) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM %s.featureofinterest WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, databaseID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}
