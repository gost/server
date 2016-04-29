package database

import (
	"fmt"
	"log"

	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/geodan/gost/sensorthings"
	"github.com/geodan/gost/sensorthings/entities"
	_ "github.com/lib/pq" // needed for PostgreSQL integration
)

// GostDatabase implementation
type GostDatabase struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Schema   string
	Ssl      bool
	Db       *sql.DB
}

// NewDatabase initialises the PostgreSQL database
//	host = TCP host:port or Unix socket depending on Network.
//	user = database user
//	password = database password
//	database = name of database
//	ssl = Whether to use secure TCP/IP connections (TLS).
func NewDatabase(host string, port int, user string, password string, database string, schema string, ssl bool) sensorthings.Database {
	return &GostDatabase{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
		Schema:   schema,
		Ssl:      ssl,
	}
}

// Start the database
func (gdb *GostDatabase) Start() {
	//ToDo: implement SSL
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", gdb.Host, gdb.User, gdb.Password, gdb.Database)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}

	gdb.Db = db
	log.Printf("Connected to database, host: \"%v\", port: \"%v\" user: \"%v\", database: \"%v\", schema: \"%v\" ssl: \"%v\"", gdb.Host, gdb.Port, gdb.User, gdb.Database, gdb.Schema, gdb.Ssl)
	//gdb.CreateSchema()
}

// CreateSchema creates the needed schema in the database
func (gdb *GostDatabase) CreateSchema() {
	create := GetCreateDatabaseQuery(gdb.Schema)
	_, err := gdb.Db.Exec(create)
	if err != nil {
		log.Fatal(err)
	}
}

// GetThing returns a thing entity based on id and query
func (gdb *GostDatabase) GetThing(id string) (*entities.Thing, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var thingID int
	var description string
	var properties string
	err2 := gdb.Db.QueryRow("select id, description, properties from v1.thing where id = $1", intID).Scan(&thingID, &description, &properties)

	if err2 != nil {
		return nil, err
	}

	thing := entities.Thing{}
	thing.ID = strconv.Itoa(thingID)
	thing.Description = description

	var p map[string]string
	err3 := json.Unmarshal([]byte(properties), &p)
	if err3 != nil {
		return nil, err3
	}

	thing.Properties = p

	return &thing, nil
}

// GetThings returns an array of things
func (gdb *GostDatabase) GetThings() ([]*entities.Thing, error) {
	rows, err := gdb.Db.Query("SELECT id, description, properties FROM v1.thing")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var things = []*entities.Thing{}

	for rows.Next() {
		thing := entities.Thing{}

		var id int
		var description string
		var properties string
		err2 := rows.Scan(&id, &description, &properties)
		if err2 != nil {
			return nil, err2
		}

		thing.ID = strconv.Itoa(id)
		thing.Description = description

		var p map[string]string
		err3 := json.Unmarshal([]byte(properties), &p)
		if err3 != nil {
			return nil, err3
		}

		thing.Properties = p
		things = append(things, &thing)
	}

	return things, nil
}

// PostThing receives a posted thing entity and adds it to the database
// returns the created Thing including the generated id
func (gdb *GostDatabase) PostThing(thing entities.Thing) (*entities.Thing, error) {
	jsonProperties, _ := json.Marshal(thing.Properties)
	var thingID int
	err := gdb.Db.QueryRow("INSERT INTO v1.thing (description, properties) VALUES ($1, $2) RETURNING id", thing.Description, jsonProperties).Scan(&thingID)
	if err != nil {
		return nil, err
	}

	thing.ID = strconv.Itoa(thingID)
	return &thing, nil
}

// PostLocation receives a posted location entity and adds it to the database
// returns the created Location including the generated id
func (gdb *GostDatabase) PostLocation(location entities.Location) (*entities.Location, error) {
	var locationID int
	err := gdb.Db.QueryRow("INSERT INTO v1.location (description, encodingtype, location) VALUES ($1, $2, $3) RETURNING id", location.Description, 1, location.Location).Scan(&locationID)
	if err != nil {
		return nil, err
	}

	location.ID = strconv.Itoa(locationID)
	return &location, nil
}

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
	_, err3 := gdb.Db.Exec("INSERT INTO v1.historicallocation (time, thing_id, location_id) VALUES ($1, $2, $3)", time.Now(), tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}

// LinkLocation links a thing with a location
// fails when a thing or location cannot be found for the given id's
func (gdb *GostDatabase) LinkLocation(thingID string, locationID string) error {
	tid, err := strconv.Atoi(thingID)
	if !gdb.ThingExists(tid) || err != nil {
		return fmt.Errorf("Thing(%v) does not exist", thingID)
	}

	lid, err2 := strconv.Atoi(thingID)
	if !gdb.ThingExists(lid) || err2 != nil {
		return fmt.Errorf("Location(%v) does not exist", locationID)
	}

	_, err3 := gdb.Db.Exec("INSERT INTO v1.thing_to_location (thing_id, location_id) VALUES ($1, $2)", tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}

// ThingExists checks if a thing is present in the database based on a given id
func (gdb *GostDatabase) ThingExists(thingID int) bool {
	var result bool
	err := gdb.Db.QueryRow("SELECT exists (SELECT 1 FROM v1.thing WHERE id = $1 LIMIT 1)", thingID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}

// LocationExists checks if a location is present in the database based on a given id
func (gdb *GostDatabase) LocationExists(locationID int) bool {
	var result bool
	err := gdb.Db.QueryRow("SELECT exists (SELECT 1 FROM v1.location WHERE id = $1 LIMIT 1)", locationID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}
