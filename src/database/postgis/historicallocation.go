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

var hlMapping = map[string]string{"time": fmt.Sprintf("to_char(time at time zone 'UTC', '%s') as time", TimeFormat)}

// GetTotalHistoricalLocations returns the amount of HistoricalLocations in the database
func (gdb *GostDatabase) GetTotalHistoricalLocations() int {
	var count int
	sql := fmt.Sprintf("SELECT Count(*) from %s.historicallocation", gdb.Schema)
	gdb.Db.QueryRow(sql).Scan(&count)
	return count
}

// GetHistoricalLocation retireves a HistoricalLocation by id
func (gdb *GostDatabase) GetHistoricalLocation(id interface{}, qo *odata.QueryOptions) (*entities.HistoricalLocation, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("HistoricalLocation does not exist"))
	}

	sql := fmt.Sprintf("select "+CreateSelectString(&entities.HistoricalLocation{}, qo, "", "", hlMapping)+" FROM %s.historicallocation where id = %v", gdb.Schema, intID)
	historicallocation, err := processHistoricalLocation(gdb.Db, sql, qo)
	if err != nil {
		return nil, err
	}

	return historicallocation, nil
}

// GetHistoricalLocations retrieves all historicallocations
func (gdb *GostDatabase) GetHistoricalLocations(qo *odata.QueryOptions) ([]*entities.HistoricalLocation, int, error) {
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.HistoricalLocation{}, qo, "", "", hlMapping)+" FROM %s.historicallocation order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.historicallocation", gdb.Schema)
	return processHistoricalLocations(gdb.Db, sql, qo, countSQL)
}

// GetHistoricalLocationsByLocation retrieves all historicallocations linked to the given location
func (gdb *GostDatabase) GetHistoricalLocationsByLocation(locationID interface{}, qo *odata.QueryOptions) ([]*entities.HistoricalLocation, int, error) {
	intID, ok := ToIntID(locationID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}
	query := fmt.Sprintf("select "+CreateSelectString(&entities.HistoricalLocation{}, qo, "", "", hlMapping)+" FROM %s.historicallocation inner join %s.location_to_historicallocation on location_to_historicallocation.historicallocation_id = historicallocation.id where location_to_historicallocation.location_id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, gdb.Schema, intID)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.historicallocation inner join %s.location_to_historicallocation on location_to_historicallocation.historicallocation_id = historicallocation.id where location_to_historicallocation.location_id = %v", gdb.Schema, gdb.Schema, intID)
	return processHistoricalLocations(gdb.Db, query, qo, countSQL)
}

// GetHistoricalLocationsByThing retrieves all historicallocations linked to the given thing
func (gdb *GostDatabase) GetHistoricalLocationsByThing(thingID interface{}, qo *odata.QueryOptions) ([]*entities.HistoricalLocation, int, error) {
	intID, ok := ToIntID(thingID)
	if !ok {
		return nil, 0, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.HistoricalLocation{}, qo, "", "", hlMapping)+" FROM %s.historicallocation where thing_id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, intID)
	countSQL := fmt.Sprintf("select COUNT(*) FROM %s.historicallocation where thing_id = %v", gdb.Schema, intID)
	return processHistoricalLocations(gdb.Db, sql, qo, countSQL)
}

func processHistoricalLocation(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.HistoricalLocation, error) {
	hls, _, err := processHistoricalLocations(db, sql, qo, "")
	if err != nil {
		return nil, err
	}

	if len(hls) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("HistoricalLocation not found"))
	}

	return hls[0], nil
}

func processHistoricalLocations(db *sql.DB, sql string, qo *odata.QueryOptions, countSQL string) ([]*entities.HistoricalLocation, int, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	var hls = []*entities.HistoricalLocation{}
	for rows.Next() {
		var id interface{}
		var time string

		var params []interface{}
		var qp []string
		if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
			s := &entities.HistoricalLocation{}
			qp = s.GetPropertyNames()
		} else {
			qp = qo.QuerySelect.Params
		}

		for _, p := range qp {
			p = strings.ToLower(p)
			if p == "id" {
				params = append(params, &id)
			}
			if p == "time" {
				params = append(params, &time)
			}
		}

		err = rows.Scan(params...)

		datastream := entities.HistoricalLocation{}
		datastream.ID = id
		datastream.Time = time

		hls = append(hls, &datastream)
	}

	var count int
	if len(countSQL) > 0 {
		db.QueryRow(countSQL).Scan(&count)
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
		query := fmt.Sprintf("INSERT INTO %s.location_to_historicallocation (location_id, historicallocation_id) VALUES ($1, $2)", gdb.Schema)
		gdb.Db.QueryRow(query, lid, hlID)
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
	/*
		var ok bool
		var intID int
		if intID, ok = ToIntID(id); !ok || !gdb.HistoricalLocationExists(intID) {
			return nil, gostErrors.NewRequestNotFound(errors.New("HistoricalLocation does not exist"))
		}

		if hl.Thing != nil {
			tid, ok := ToIntID(hl.Thing.ID)
			if !ok || !gdb.ThingExists(tid) {
				return nil, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
			}

			query := fmt.Sprintf("UPDATE %s.historicallocation set thing_id=$1 where id = $2", gdb.Schema)
			_, err := gdb.Db.Exec(query, intID)
			if err != nil {
				return nil, err
			}
		}

		if len(hl.Locations) > 0 {
			for _, l := range hl.Locations {
				lid, ok := ToIntID(l.ID)
				if !ok || !gdb.LocationExists(lid) {
					return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
				}
			}

			query := fmt.Sprintf("DELETE FROM %s.location_to_historicallocation WHERE historicallocation_id = $1", gdb.Schema)
			_, err := gdb.Db.Exec(query, intID)
			if err != nil {
				return nil, err
			}

			for _, l := range hl.Locations {
				query = fmt.Sprintf("INSERT INTO %s.location_to_historicallocation (location_id, historicallocation_id) VALUES ($1, $2)", gdb.Schema)
				_, err = gdb.Db.Exec(query, l.ID, intID)
				if err != nil {
					return nil, err
				}
			}
		}

		if len(hl.Time) > 0 {
			t, _ := time.Parse(time.RFC3339Nano, hl.Time)
			utcT := t.UTC().Format("2006-01-02T15:04:05.000Z")

			query := fmt.Sprintf("UPDATE %s.historicallocation set time=$1 where id = $2", gdb.Schema)
			_, err := gdb.Db.Exec(query, utcT, intID)
			if err != nil {
				return nil, err
			}
		}

		hl.ID = intID

		return hl, nil
	*/
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

	/*
		query := fmt.Sprintf("DELETE FROM %s.location_to_historicallocation WHERE historicallocation_id = $1", gdb.Schema)
		_, err = gdb.Db.Exec(query, intID)
		if err != nil {
			return nil, err
		}*/

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
