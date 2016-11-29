package postgis

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
	"strings"
)

var foiMapping = map[string]string{"feature": "public.ST_AsGeoJSON(featureofinterest.feature) AS feature"}

// GetFeatureOfInterestByLocationID returns the FeatureOfInterest in the database
// where original_location_id equals the given parameter
func (gdb *GostDatabase) GetFeatureOfInterestByLocationID(id interface{}) (*entities.FeatureOfInterest, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.FeatureOfInterest{}, nil, "", "", foiMapping)+" from %s.featureofinterest where original_location_id=%v", gdb.Schema, intID)
	return processFeatureOfInterest(gdb.Db, sql, nil)
}

// GetTotalFeaturesOfInterest returns the number of FeaturesOfInterest records in the database
func (gdb *GostDatabase) GetTotalFeaturesOfInterest() int {
	var count int
	sql := fmt.Sprintf("SELECT Count(*) from %s.featureofinterest", gdb.Schema)
	gdb.Db.QueryRow(sql).Scan(&count)
	return count
}

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
		return nil, gostErrors.NewRequestNotFound(errors.New("Observation does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.FeatureOfInterest{}, qo, "featureofinterest.", "", foiMapping)+" from %s.featureofinterest inner join %s.observation on observation.featureofinterest_id = featureofinterest.id where observation.id = %v limit 1", gdb.Schema, gdb.Schema, intID)
	return processFeatureOfInterest(gdb.Db, sql, qo)
}

// GetFeatureOfInterests returns all feature of interests
func (gdb *GostDatabase) GetFeatureOfInterests(qo *odata.QueryOptions) ([]*entities.FeatureOfInterest, int, error) {
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.FeatureOfInterest{}, qo, "", "", foiMapping)+" from %s.featureofinterest order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema)
	countSSQL := fmt.Sprintf("select COUNT(*) FROM %s.featureofinterest", gdb.Schema)
	return processFeatureOfInterests(gdb.Db, sql, qo, countSSQL)
}

// PostFeatureOfInterest inserts a new FeatureOfInterest into the database
func (gdb *GostDatabase) PostFeatureOfInterest(f *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	var fID int
	locationBytes, _ := json.Marshal(f.Feature)
	encoding, _ := entities.CreateEncodingType(f.EncodingType)
	sql := fmt.Sprintf("INSERT INTO %s.featureofinterest (name, description, encodingtype, feature, original_location_id) VALUES ($1, $2, $3, ST_SetSRID(public.ST_GeomFromGeoJSON('%s'),4326), $4) RETURNING id", gdb.Schema, string(locationBytes[:]))
	err := gdb.Db.QueryRow(sql, f.Name, f.Description, encoding.Code, f.OriginalLocationID).Scan(&fID)
	if err != nil {
		return nil, err
	}

	f.ID = fID
	return f, nil
}

// PutFeatureOfInterest inserts a new FeatureOfInterest into the database
func (gdb *GostDatabase) PutFeatureOfInterest(id interface{}, f *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	return gdb.PatchFeatureOfInterest(id, f)
	/*
		locationBytes, _ := json.Marshal(f.Feature)
		intID, _ := ToIntID(id)
		encoding, _ := entities.CreateEncodingType(f.EncodingType)

		if !gdb.FeatureOfInterestExists(intID) {
			return nil, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest does not exist"))
		}

		sql := fmt.Sprintf("update %s.featureofinterest set name=$1, description=$2, encodingtype=$3, feature= ST_SetSRID(public.ST_GeomFromGeoJSON('%s'),4326), original_location_id=$4 where id=$5", gdb.Schema, string(locationBytes[:]))
		_, err := gdb.Db.Exec(sql, f.Name, f.Description, encoding.Code, f.OriginalLocationID, intID)
		if err != nil {
			return nil, err
		}

		f.ID = intID
		return f, nil*/
}

func processFeatureOfInterest(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error) {
	locations, _, err := processFeatureOfInterests(db, sql, qo, "")
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest not found"))
	}

	return locations[0], nil
}

func processFeatureOfInterests(db *sql.DB, sql string, qo *odata.QueryOptions, countSQL string) ([]*entities.FeatureOfInterest, int, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	var featureOfInterests = []*entities.FeatureOfInterest{}
	for rows.Next() {
		var ID interface{}
		var encodingType int
		var name, description, feature string

		var params []interface{}
		var qp []string
		if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
			f := &entities.FeatureOfInterest{}
			qp = f.GetPropertyNames()
		} else {
			qp = qo.QuerySelect.Params
		}

		for _, p := range qp {
			p = strings.ToLower(p)
			if p == "id" {
				params = append(params, &ID)
			}
			if p == "encodingtype" {
				params = append(params, &encodingType)
			}
			if p == "name" {
				params = append(params, &name)
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
			return nil, 0, err
		}

		foi := entities.FeatureOfInterest{}
		foi.ID = ID
		foi.Name = name
		foi.Description = description
		foi.Feature = featureMap
		if encodingType != 0 {
			foi.EncodingType = entities.EncodingValues[encodingType].Value
		}

		featureOfInterests = append(featureOfInterests, &foi)
	}

	var count int
	if len(countSQL) > 0 {
		db.QueryRow(countSQL).Scan(&count)
	}

	return featureOfInterests, count, nil
}

// PatchFeatureOfInterest updates a FeatureOfInterest in the database
func (gdb *GostDatabase) PatchFeatureOfInterest(id interface{}, foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	var err error
	var ok bool
	var intID int
	updates := make(map[string]interface{})

	if intID, ok = ToIntID(id); !ok || !gdb.FeatureOfInterestExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest does not exist"))
	}

	if len(foi.Name) > 0 {
		updates["name"] = foi.Name
	}

	if len(foi.Description) > 0 {
		updates["description"] = foi.Description
	}

	if len(foi.EncodingType) > 0 {
		encoding, _ := entities.CreateEncodingType(foi.EncodingType)
		updates["encodingtype"] = encoding.Code
	}

	if len(foi.Feature) > 0 {
		locationBytes, _ := json.Marshal(foi.Feature)
		updates["feature"] = fmt.Sprintf("ST_SetSRID(public.ST_GeomFromGeoJSON('%s'),4326)", string(locationBytes[:]))
	}

	if err = gdb.updateEntityColumns("featureofinterest", updates, intID); err != nil {
		return nil, err
	}

	nfoi, _ := gdb.GetFeatureOfInterest(intID, nil)
	return nfoi, nil
}

// DeleteFeatureOfInterest tries to delete a FeatureOfInterest by the given id
func (gdb *GostDatabase) DeleteFeatureOfInterest(id interface{}) error {
	return DeleteEntity(gdb, id, "featureofinterest")
}

// FeatureOfInterestExists checks if a FeatureOfInterest is present in the database based on a given id.
func (gdb *GostDatabase) FeatureOfInterestExists(id int) bool {
	return EntityExists(gdb, id, "featureofinterest")
}
