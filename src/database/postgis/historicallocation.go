package postgis

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
	"strings"
)

func historicalLocationParamFactory(values map[string]interface{}) (entities.Entity, error) {
	h := &entities.HistoricalLocation{}
	for as, value := range values {
		if value == nil {
			continue
		}

		if as == asMappings[entities.EntityTypeHistoricalLocation][historicalLocationID] {
			h.ID = value
		} else if as == asMappings[entities.EntityTypeHistoricalLocation][historicalLocationTime] {
			t := value.(string)
			h.Time = strings.Replace(t, "\"", "", 2)
		}
	}

	return h, nil
}

// GetHistoricalLocation retrieves a HistoricalLocation by id
func (gdb *GostDatabase) GetHistoricalLocation(id interface{}, qo *odata.QueryOptions) (*entities.HistoricalLocation, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("HistoricalLocation does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.HistoricalLocation{}, nil, intID, qo)
	return processHistoricalLocation(gdb.Db, query, qi)
}

// GetHistoricalLocations retrieves all historicallocations
func (gdb *GostDatabase) GetHistoricalLocations(qo *odata.QueryOptions) ([]*entities.HistoricalLocation, int, error) {
	query, qi := gdb.QueryBuilder.CreateQuery(&entities.HistoricalLocation{}, nil, nil, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.HistoricalLocation{}, nil, nil, qo)
	return processHistoricalLocations(gdb.Db, query, qi, countSQL)
}

// GetHistoricalLocationsByLocation retrieves all historicallocations linked to the given location
func (gdb *GostDatabase) GetHistoricalLocationsByLocation(locationID interface{}, qo *odata.QueryOptions) ([]*entities.HistoricalLocation, int, error) {
	intID, ok := ToIntID(locationID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}
	query, qi := gdb.QueryBuilder.CreateQuery(&entities.HistoricalLocation{}, &entities.Location{}, intID, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.HistoricalLocation{}, &entities.Location{}, intID, qo)
	return processHistoricalLocations(gdb.Db, query, qi, countSQL)
}

// GetHistoricalLocationsByThing retrieves all historicallocations linked to the given thing
func (gdb *GostDatabase) GetHistoricalLocationsByThing(thingID interface{}, qo *odata.QueryOptions) ([]*entities.HistoricalLocation, int, error) {
	intID, ok := ToIntID(thingID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}
	query, qi := gdb.QueryBuilder.CreateQuery(&entities.HistoricalLocation{}, &entities.Thing{}, intID, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.HistoricalLocation{}, &entities.Thing{}, intID, qo)
	return processHistoricalLocations(gdb.Db, query, qi, countSQL)
}

func processHistoricalLocation(db *sql.DB, sql string, qi *QueryParseInfo) (*entities.HistoricalLocation, error) {
	hls, _, err := processHistoricalLocations(db, sql, qi, "")
	if err != nil {
		return nil, err
	}

	if len(hls) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("HistoricalLocation not found"))
	}

	return hls[0], nil
}

func processHistoricalLocations(db *sql.DB, sql string, qi *QueryParseInfo, countSQL string) ([]*entities.HistoricalLocation, int, error) {
	data, err := ExecuteSelect(db, qi, sql)
	if err != nil {
		return nil, 0, fmt.Errorf("Error executing query %v", err)
	}

	hls := make([]*entities.HistoricalLocation, 0)
	for _, d := range data {
		entity := d.(*entities.HistoricalLocation)
		hls = append(hls, entity)
	}

	var count int
	if len(countSQL) > 0 {
		count, err = ExecuteSelectCount(db, countSQL)
		if err != nil {
			return nil, 0, fmt.Errorf("Error executing count %v", err)
		}
	}

	return hls, count, nil
}

// PostHistoricalLocation adds a historical location to the database
// returns the created historical location including the generated id
// fails when a thing or location cannot be found for the given id's
func (gdb *GostDatabase) PostHistoricalLocation(hl *entities.HistoricalLocation) (*entities.HistoricalLocation, error) {
	var hlID int
	var err error
	tid, ok := ToIntID(hl.Thing.ID)
	if !ok || !gdb.ThingExists(tid) {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}

	for _, l := range hl.Locations {
		lid, ok := ToIntID(l.ID)
		if !ok || !gdb.LocationExists(lid) {
			return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
		}
	}

	query := fmt.Sprintf("INSERT INTO %s.historicallocation (time, thing_id) VALUES ($1, $2) RETURNING id", gdb.Schema)
	err = gdb.Db.QueryRow(query, time.Now(), tid).Scan(&hlID)
	if err != nil {
		return nil, err
	}

	for _, l := range hl.Locations {
		lid, _ := ToIntID(l.ID)
		query := fmt.Sprintf("INSERT INTO %s.location_to_historicallocation (location_id, historicallocation_id) VALUES ($1, $2)  RETURNING historicallocation_id", gdb.Schema)
		err = gdb.Db.QueryRow(query, lid, hlID).Scan(&lid)
		if err != nil {
			return nil, err
		}
	}

	hl.ID = hlID
	hl.Locations = nil
	return hl, nil
}

// HistoricalLocationExists checks if a HistoricalLocation is present in the database based on a given id
func (gdb *GostDatabase) HistoricalLocationExists(id interface{}) bool {
	return EntityExists(gdb, id, "historicallocation")
}

// PutHistoricalLocation updates a HistoricalLocation in the database
func (gdb *GostDatabase) PutHistoricalLocation(id interface{}, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, error) {
	return gdb.PatchHistoricalLocation(id, hl)
}

// PatchHistoricalLocation updates a HistoricalLocation in the database
func (gdb *GostDatabase) PatchHistoricalLocation(id interface{}, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, error) {
	var err error
	var ok bool
	var intID int
	updates := make(map[string]interface{})

	if intID, ok = ToIntID(id); !ok || !gdb.HistoricalLocationExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("HistoricalLocation does not exist"))
	}

	for _, l := range hl.Locations {
		lid, ok := ToIntID(l.ID)
		if !ok || !gdb.LocationExists(lid) {
			return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
		}
	}

	if len(hl.Time) > 0 {
		updates["time"] = hl.Time
	}

	if err = gdb.updateEntityColumns("historicallocation", updates, intID); err != nil {
		return nil, err
	}

	for _, l := range hl.Locations {
		query := fmt.Sprintf("INSERT INTO %s.location_to_historicallocation (location_id, historicallocation_id) VALUES ($1, $2)", gdb.Schema)
		_, err := gdb.Db.Exec(query, l.ID, intID)
		if err != nil {
			return nil, err
		}
	}

	nhl, _ := gdb.GetHistoricalLocation(intID, nil)
	return nhl, nil
}

// DeleteHistoricalLocation tries to delete a HistoricalLocation by the given id
func (gdb *GostDatabase) DeleteHistoricalLocation(id interface{}) error {
	return DeleteEntity(gdb, id, "historicallocation")
}
