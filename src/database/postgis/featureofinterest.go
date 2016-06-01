package postgis

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"strconv"
)

// GetFeatureOfInterest returns a feature of interest by id
func (gdb *GostDatabase) GetFeatureOfInterest(id string) (*entities.FeatureOfInterest, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select id, description, encodingtype, public.ST_AsGeoJSON(feature) AS feature from %s.featureofinterest where id = $1", gdb.Schema)
	return processFeatureOfInterest(gdb.Db, sql, intID)
}

// GetFeatureOfInterestByObservation returns a feature of interest by given observation id
func (gdb *GostDatabase) GetFeatureOfInterestByObservation(id string) (*entities.FeatureOfInterest, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select featureofinterest.id, featureofinterest.description, featureofinterest.encodingtype, public.ST_AsGeoJSON(featureofinterest.feature) AS feature from %s.featureofinterest inner join %s.observation on observation.featureofinterest_id = featureofinterest.id where observation.id = $1 limit 1", gdb.Schema, gdb.Schema)
	return processFeatureOfInterest(gdb.Db, sql, intID)
}

// GetFeatureOfInterests returns all feature of interests
func (gdb *GostDatabase) GetFeatureOfInterests() ([]*entities.FeatureOfInterest, error) {
	sql := fmt.Sprintf("select id, description, encodingtype, public.ST_AsGeoJSON(feature) AS feature from %s.featureofinterest", gdb.Schema)
	return processFeatureOfInterests(gdb.Db, sql)
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

func processFeatureOfInterest(db *sql.DB, sql string, args ...interface{}) (*entities.FeatureOfInterest, error) {
	locations, err := processFeatureOfInterests(db, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest not found"))
	}

	return locations[0], nil
}

func processFeatureOfInterests(db *sql.DB, sql string, args ...interface{}) ([]*entities.FeatureOfInterest, error) {
	rows, err := db.Query(sql, args...)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var featureOfInterests = []*entities.FeatureOfInterest{}
	for rows.Next() {
		var ID, encodingtype int
		var description, feature string
		err = rows.Scan(&ID, &description, &encodingtype, &feature)
		if err != nil {
			return nil, err
		}

		featureMap, err := JSONToMap(&feature)
		if err != nil {
			return nil, err
		}

		foi := entities.FeatureOfInterest{}
		foi.ID = strconv.Itoa(ID)
		foi.Description = description
		foi.Feature = featureMap
		foi.EncodingType = entities.EncodingValues[encodingtype].Value

		featureOfInterests = append(featureOfInterests, &foi)
	}

	return featureOfInterests, nil
}

// DeleteFeatureOfInterest tries to delete a FeatureOfInterest by the given id
func (gdb *GostDatabase) DeleteFeatureOfInterest(id string) error {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
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
