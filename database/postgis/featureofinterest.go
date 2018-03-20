package postgis

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	entities "github.com/gost/core"
	gostErrors "github.com/gost/server/errors"
	"github.com/gost/server/sensorthings/odata"
)

func featureOfInterestParamFactory(values map[string]interface{}) (entities.Entity, error) {
	foi := &entities.FeatureOfInterest{}
	for as, value := range values {
		if value == nil {
			continue
		}

		if as == asMappings[entities.EntityTypeFeatureOfInterest][foiID] {
			foi.ID = value
		} else if as == asMappings[entities.EntityTypeFeatureOfInterest][foiName] {
			foi.Name = value.(string)
		} else if as == asMappings[entities.EntityTypeFeatureOfInterest][foiDescription] {
			foi.Description = value.(string)
		} else if as == asMappings[entities.EntityTypeFeatureOfInterest][foiEncodingType] {
			encodingType := value.(int64)
			if encodingType != 0 {
				foi.EncodingType = entities.EncodingValues[encodingType].Value
			}
		} else if as == asMappings[entities.EntityTypeFeatureOfInterest][foiFeature] || as == asMappings[entities.EntityTypeFeatureOfInterest][foiGeoJSON] {
			t := value.(string)
			featureMap, err := JSONToMap(&t)
			if err != nil {
				return nil, err
			}

			foi.Feature = featureMap
		}
	}

	return foi, nil
}

// GetFeatureOfInterestIDByLocationID returns the FeatureOfInterest id in the database
// where original_location_id equals the given parameter
func (gdb *GostDatabase) GetFeatureOfInterestIDByLocationID(id interface{}) (interface{}, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	var fID interface{}
	query := fmt.Sprintf("select id from %s.featureofinterest where original_location_id=%v", gdb.Schema, intID)
	err := gdb.Db.QueryRow(query).Scan(&fID)
	if err != nil {
		return nil, err
	}

	if fID == nil {
		return nil, errors.New("Linked FeatureOfINterest not found")
	}

	return fID, nil
}

// GetFeatureOfInterest returns a feature of interest by id
func (gdb *GostDatabase) GetFeatureOfInterest(id interface{}, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.FeatureOfInterest{}, nil, intID, qo)
	return processFeatureOfInterest(gdb.Db, query, qi)
}

// GetFeatureOfInterestByObservation returns a feature of interest by given observation id
func (gdb *GostDatabase) GetFeatureOfInterestByObservation(id interface{}, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Observation does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.FeatureOfInterest{}, &entities.Observation{}, intID, qo)
	return processFeatureOfInterest(gdb.Db, query, qi)
}

// GetFeatureOfInterests returns all feature of interests
func (gdb *GostDatabase) GetFeatureOfInterests(qo *odata.QueryOptions) ([]*entities.FeatureOfInterest, int, bool, error) {
	query, qi := gdb.QueryBuilder.CreateQuery(&entities.FeatureOfInterest{}, nil, nil, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.FeatureOfInterest{}, nil, nil, qo)
	return processFeatureOfInterests(gdb.Db, query, qo, qi, countSQL)
}

// PostFeatureOfInterest inserts a new FeatureOfInterest into the database
func (gdb *GostDatabase) PostFeatureOfInterest(f *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	var fID int
	locationBytes, _ := json.Marshal(f.Feature)
	encoding, _ := entities.CreateEncodingType(f.EncodingType)
	sql2 := fmt.Sprintf("INSERT INTO %s.featureofinterest (name, description, encodingtype, feature, original_location_id, geojson) VALUES ($1, $2, $3, ST_SetSRID(public.ST_GeomFromGeoJSON('%s'),4326), $4, $5) RETURNING id", gdb.Schema, string(locationBytes[:]))
	err := gdb.Db.QueryRow(sql2, f.Name, f.Description, encoding.Code, f.OriginalLocationID, string(locationBytes[:])).Scan(&fID)
	if err != nil {
		return nil, err
	}

	f.ID = fID
	return f, nil
}

// PutFeatureOfInterest inserts a new FeatureOfInterest into the database
func (gdb *GostDatabase) PutFeatureOfInterest(id interface{}, f *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	return gdb.PatchFeatureOfInterest(id, f)
}

func processFeatureOfInterest(db *sql.DB, sql string, qi *QueryParseInfo) (*entities.FeatureOfInterest, error) {
	locations, _, _, err := processFeatureOfInterests(db, sql, nil, qi, "")
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("FeatureOfInterest not found"))
	}

	return locations[0], nil
}

func processFeatureOfInterests(db *sql.DB, sql string, qo *odata.QueryOptions, qi *QueryParseInfo, countSQL string) ([]*entities.FeatureOfInterest, int, bool, error) {
	data, hasNext, err := ExecuteSelect(db, qi, sql, qo)
	if err != nil {
		return nil, 0, false, fmt.Errorf("Error executing query %v", err)
	}

	fois := make([]*entities.FeatureOfInterest, 0)
	for _, d := range data {
		entity := d.(*entities.FeatureOfInterest)
		fois = append(fois, entity)
	}

	var count int
	if len(countSQL) > 0 {
		count, err = ExecuteSelectCount(db, countSQL)
		if err != nil {
			return nil, 0, false, fmt.Errorf("Error executing count %v", err)
		}
	}

	return fois, count, hasNext, nil
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
		updates["geojson"] = string(locationBytes[:])
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
