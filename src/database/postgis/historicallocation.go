package postgis

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetHistoricalLocation retireves a HistoricalLocation by id
func (gdb *GostDatabase) GetHistoricalLocation(id string, qo *odata.QueryOptions) (*entities.HistoricalLocation, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select id, to_char(time at time zone 'UTC', '%s') as time FROM %s.historicallocation where id = $1", TimeFormat, gdb.Schema)
	historicallocation, err := processHistoricalLocation(gdb.Db, sql, intID)
	if err != nil {
		return nil, err
	}

	return historicallocation, nil
}

// GetHistoricalLocations retrieves all historicallocations
func (gdb *GostDatabase) GetHistoricalLocations(qo *odata.QueryOptions) ([]*entities.HistoricalLocation, error) {
	sql := fmt.Sprintf("select id, to_char(time at time zone 'UTC', '%s') as time FROM %s.historicallocation", TimeFormat, gdb.Schema)
	return processHistoricalLocations(gdb.Db, sql)
}

// GetHistoricalLocationsByLocation retrieves all historicallocations linked to the given location
func (gdb *GostDatabase) GetHistoricalLocationsByLocation(thingID string, qo *odata.QueryOptions) ([]*entities.HistoricalLocation, error) {
	tID, err := strconv.Atoi(thingID)
	if err != nil {
		return nil, err
	}
	sql := fmt.Sprintf("select id, to_char(time at time zone 'UTC', '%s') as time FROM %s.historicallocation where location_id = $1", TimeFormat, gdb.Schema)
	return processHistoricalLocations(gdb.Db, sql, tID)
}

// GetHistoricalLocationsByThing retrieves all historicallocations linked to the given thing
func (gdb *GostDatabase) GetHistoricalLocationsByThing(thingID string, qo *odata.QueryOptions) ([]*entities.HistoricalLocation, error) {
	tID, err := strconv.Atoi(thingID)
	if err != nil {
		return nil, err
	}
	sql := fmt.Sprintf("select id, to_char(time at time zone 'UTC', '%s') as time FROM %s.historicallocation where thing_id = $1", TimeFormat, gdb.Schema)
	return processHistoricalLocations(gdb.Db, sql, tID)
}

func processHistoricalLocation(db *sql.DB, sql string, args ...interface{}) (*entities.HistoricalLocation, error) {
	hls, err := processHistoricalLocations(db, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(hls) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("HistoricalLocation not found"))
	}

	return hls[0], nil
}

func processHistoricalLocations(db *sql.DB, sql string, args ...interface{}) ([]*entities.HistoricalLocation, error) {
	rows, err := db.Query(sql, args...)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var hls = []*entities.HistoricalLocation{}
	for rows.Next() {
		var id int
		var time string

		err := rows.Scan(&id, &time)
		if err != nil {
			return nil, err
		}

		datastream := entities.HistoricalLocation{}
		datastream.ID = strconv.Itoa(id)
		datastream.Time = time

		hls = append(hls, &datastream)
	}

	return hls, nil
}

// PostHistoricalLocation adds a historical location to the database
// returns the created historical location including the generated id
// fails when a thing or location cannot be found for the given id's
func (gdb *GostDatabase) PostHistoricalLocation(thingID string, locationID string) error {
	tid, err := strconv.Atoi(thingID)
	if !gdb.ThingExists(tid) || err != nil {
		return fmt.Errorf("Thing(%s) does not exist", thingID)
	}

	lid, err2 := strconv.Atoi(locationID)
	if !gdb.LocationExists(lid) || err2 != nil {
		return fmt.Errorf("Location(%s) does not exist", locationID)
	}

	sql := fmt.Sprintf("INSERT INTO %s.historicallocation (time, thing_id, location_id) VALUES ($1, $2, $3)", gdb.Schema)
	_, err3 := gdb.Db.Exec(sql, time.Now(), tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}

// DeleteHistoricalLocation tries to delete a HistoricalLocation by the given id
func (gdb *GostDatabase) DeleteHistoricalLocation(id string) error {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
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
