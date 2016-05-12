package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
)

// GetFeatureOfInterest todo
func (gdb *GostDatabase) GetFeatureOfInterest(id string) (*entities.FeatureOfInterest, error) {
	return nil, nil
}

// GetFeatureOfInterests todo
func (gdb *GostDatabase) GetFeatureOfInterests() ([]*entities.FeatureOfInterest, error) {
	return nil, nil
}

// PostFeatureOfInterest todo
func (gdb *GostDatabase) PostFeatureOfInterest(f entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	return nil, nil
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
