package postgis

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
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
func (gdb *GostDatabase) GetHistoricalLocations(qo *odata.QueryOptions) ([]*entities.HistoricalLocation, error) {
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.HistoricalLocation{}, qo, "", "", hlMapping)+" FROM %s.historicallocation order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema)
	return processHistoricalLocations(gdb.Db, sql, qo)
}

// GetHistoricalLocationsByLocation retrieves all historicallocations linked to the given location
func (gdb *GostDatabase) GetHistoricalLocationsByLocation(locationID interface{}, qo *odata.QueryOptions) ([]*entities.HistoricalLocation, error) {
	tID, ok := ToIntID(locationID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.HistoricalLocation{}, qo, "", "", hlMapping)+" FROM %s.historicallocation where location_id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, tID)
	return processHistoricalLocations(gdb.Db, sql, qo)
}

// GetHistoricalLocationsByThing retrieves all historicallocations linked to the given thing
func (gdb *GostDatabase) GetHistoricalLocationsByThing(thingID interface{}, qo *odata.QueryOptions) ([]*entities.HistoricalLocation, error) {
	tID, ok := ToIntID(thingID)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}
	sql := fmt.Sprintf("select "+CreateSelectString(&entities.HistoricalLocation{}, qo, "", "", hlMapping)+" FROM %s.historicallocation where thing_id = %v order by id desc "+CreateTopSkipQueryString(qo), gdb.Schema, tID)
	return processHistoricalLocations(gdb.Db, sql, qo)
}

func processHistoricalLocation(db *sql.DB, sql string, qo *odata.QueryOptions) (*entities.HistoricalLocation, error) {
	hls, err := processHistoricalLocations(db, sql, qo)
	if err != nil {
		return nil, err
	}

	if len(hls) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("HistoricalLocation not found"))
	}

	return hls[0], nil
}

func processHistoricalLocations(db *sql.DB, sql string, qo *odata.QueryOptions) ([]*entities.HistoricalLocation, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, err
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

	return hls, nil
}

// PostHistoricalLocation adds a historical location to the database
// returns the created historical location including the generated id
// fails when a thing or location cannot be found for the given id's
func (gdb *GostDatabase) PostHistoricalLocation(thingID interface{}, locationID interface{}) error {
	tid, ok := ToIntID(thingID)
	if !ok || !gdb.ThingExists(tid) {
		return gostErrors.NewRequestNotFound(errors.New("Thing does not exist"))
	}

	lid, ok := ToIntID(locationID)
	if !ok || !gdb.LocationExists(lid) {
		return gostErrors.NewRequestNotFound(errors.New("Location does not exist"))
	}

	sql := fmt.Sprintf("INSERT INTO %s.historicallocation (time, thing_id, location_id) VALUES ($1, $2, $3)", gdb.Schema)
	_, err3 := gdb.Db.Exec(sql, time.Now(), tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}

// DeleteHistoricalLocation tries to delete a HistoricalLocation by the given id
func (gdb *GostDatabase) DeleteHistoricalLocation(id interface{}) error {
	intID, ok := ToIntID(id)
	if !ok {
		return gostErrors.NewRequestNotFound(errors.New("HistoricalLocation does not exist"))
	}

	r, err := gdb.Db.Exec(fmt.Sprintf("DELETE FROM %s.historicallocation WHERE id = $1", gdb.Schema), intID)
	if err != nil {
		return err
	}

	if c, _ := r.RowsAffected(); c == 0 {
		return gostErrors.NewRequestNotFound(errors.New("HistoricalLocation not found"))
	}

	return nil
}
