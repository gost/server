package postgis

import (
	"fmt"
	"strconv"
	"time"
)

// PostHistoricalLocation adds a historical location to the database
// returns the created historical location including the generated id
// fails when a thing or location cannot be found for the given id's
func (gdb *GostDatabase) PostHistoricalLocation(thingID string, locationID string) error {
	tid, err := strconv.Atoi(thingID)
	if !gdb.ThingExists(tid) || err != nil {
		return fmt.Errorf("Thing(%v) does not exist", thingID)
	}

	lid, err2 := strconv.Atoi(thingID)
	if !gdb.ThingExists(lid) || err2 != nil {
		return fmt.Errorf("Location(%v) does not exist", locationID)
	}

	//check if thing and location exist
	sql := fmt.Sprintf("INSERT INTO %s.historicallocation (time, thing_id, location_id) VALUES ($1, $2, $3)", gdb.Schema)
	_, err3 := gdb.Db.Exec(sql, time.Now(), tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}
