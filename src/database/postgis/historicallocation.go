package postgis

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
)

// GetHistoricalLocation retireves a HistoricalLocation by id
func (gdb *GostDatabase) GetHistoricalLocation(id string) (*entities.HistoricalLocation, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("select id, to_char(time at time zone 'UTC', '%s') as time FROM historicallocation where id = $1", TimeFormat)
	historicallocation, err := processHistoricalLocation(gdb.Db, sql, intID)
	if err != nil {
		return nil, err
	}

	return historicallocation, nil
}

// GetHistoricalLocations retrieves all historicallocations
func (gdb *GostDatabase) GetHistoricalLocations() ([]*entities.HistoricalLocation, error) {
	sql := fmt.Sprintf("select id, to_char(time at time zone 'UTC', '%s') as time FROM historicallocation", TimeFormat)
	return processHistoricalLocations(gdb.Db, sql)
}

// GetHistoricalLocationsByThing retrieves all historicallocations linked to the given thing
func (gdb *GostDatabase) GetHistoricalLocationsByThing(thingID string) ([]*entities.HistoricalLocation, error) {
	tID, err := strconv.Atoi(thingID)
	if err != nil {
		return nil, err
	}
	sql := fmt.Sprintf("select id, to_char(time at time zone 'UTC', '%s') as time FROM historicallocation where thing_id = $1", TimeFormat)
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

		datastream := entities.HistoricalLocation{
			ID:   strconv.Itoa(id),
			Time: time,
		}
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

	lid, err2 := strconv.Atoi(thingID)
	if !gdb.ThingExists(lid) || err2 != nil {
		return fmt.Errorf("Location(%s) does not exist", locationID)
	}

	sql := "INSERT INTO historicallocation (time, thing_id, location_id) VALUES ($1, $2, $3)"
	_, err3 := gdb.Db.Exec(sql, time.Now(), tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}
