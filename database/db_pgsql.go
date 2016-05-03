package database

import (
	"fmt"
	"log"

	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	_ "github.com/lib/pq" // postgres driver
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
func NewDatabase(host string, port int, user string, password string, database string, schema string, ssl bool) models.Database {
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
	gdb.CreateSchema()
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
	sql := fmt.Sprintf("select id, description, properties from %s.thing where id = $1", gdb.Schema)
	err2 := gdb.Db.QueryRow(sql, intID).Scan(&thingID, &description, &properties)

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
	sql := fmt.Sprintf("select id, description, properties FROM %s.thing", gdb.Schema)
	rows, err := gdb.Db.Query(sql)
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
	sql := fmt.Sprintf("INSERT INTO %s.thing (description, properties) VALUES ($1, $2) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thing.Description, jsonProperties).Scan(&thingID)
	if err != nil {
		return nil, err
	}

	thing.ID = strconv.Itoa(thingID)
	return &thing, nil
}

// GetLocation retrieves the location for the given id from the database
func (gdb *GostDatabase) GetLocation(id string) (*entities.Location, error) {
	//example error @bertt when location cannot be found
	err := gostErrors.NewRequestNotFound(fmt.Errorf("Location(%s) does not exist", id))
	return nil, err
}

// GetLocations todo
func (gdb *GostDatabase) GetLocations() ([]*entities.Location, error) {
	return nil, nil
}

// PostLocation receives a posted location entity and adds it to the database
// returns the created Location including the generated id
func (gdb *GostDatabase) PostLocation(location entities.Location) (*entities.Location, error) {
	var locationID int
	sql := fmt.Sprintf("INSERT INTO %s.location (description, encodingtype, location) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, location.Description, 1, location.Location).Scan(&locationID)
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
	sql := fmt.Sprintf("INSERT INTO %s.historicallocation (time, thing_id, location_id) VALUES ($1, $2, $3)", gdb.Schema)
	_, err3 := gdb.Db.Exec(sql, time.Now(), tid, lid)
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

	sql := fmt.Sprintf("INSERT INTO %s.thing_to_location (thing_id, location_id) VALUES ($1, $2)", gdb.Schema)
	_, err3 := gdb.Db.Exec(sql, tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}

// GetObservedProperty todo
func (gdb *GostDatabase) GetObservedProperty(id string) (*entities.ObservedProperty, error) {
	return nil, nil
}

// GetObservedProperties todo
func (gdb *GostDatabase) GetObservedProperties() ([]*entities.ObservedProperty, error) {
	return nil, nil
}

// PostObservedProperty todo
func (gdb *GostDatabase) PostObservedProperty(op entities.ObservedProperty) (*entities.ObservedProperty, error) {
	return nil, nil
}

// GetSensor todo
func (gdb *GostDatabase) GetSensor(id string) (*entities.Sensor, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var sensorID int
	var description string
	var metadata string
	sql := fmt.Sprintf("select id, description, metadata from %s.sensor where id = $1", gdb.Schema)
	err2 := gdb.Db.QueryRow(sql, intID).Scan(&sensorID, &description, &metadata)

	if err2 != nil {
		return nil, err
	}

	sensor := entities.Sensor{}
	sensor.ID = strconv.Itoa(sensorID)
	sensor.Description = description
	sensor.Metadata = metadata

	return &sensor, nil
}

// GetSensors todo
func (gdb *GostDatabase) GetSensors() ([]*entities.Sensor, error) {
	sql := fmt.Sprintf("select id, description, metadata FROM %s.sensor", gdb.Schema)
	rows, err := gdb.Db.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var sensors = []*entities.Sensor{}

	for rows.Next() {
		sensor := entities.Sensor{}

		var id int
		var description string
		var metadata string
		err2 := rows.Scan(&id, &description, &metadata)
		if err2 != nil {
			return nil, err2
		}

		sensor.ID = strconv.Itoa(id)
		sensor.Description = description
		sensor.Metadata = metadata

		sensors = append(sensors, &sensor)
	}

	return sensors, nil
}

// PostSensor todo
func (gdb *GostDatabase) PostSensor(sensor entities.Sensor) (*entities.Sensor, error) {
	var sensorID int
	sql := fmt.Sprintf("INSERT INTO %s.sensor (description, encodingtype, metadata) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)
	err := gdb.Db.QueryRow(sql, sensor.Description, 1, sensor.Metadata).Scan(&sensorID)
	if err != nil {
		return nil, err
	}

	sensor.ID = strconv.Itoa(sensorID)
	return &sensor, nil
}

// GetDatastream todo
func (gdb *GostDatabase) GetDatastream(id string) (*entities.Datastream, error) {
	return nil, nil
}

// GetDatastreams todo
func (gdb *GostDatabase) GetDatastreams() ([]*entities.Datastream, error) {
	return nil, nil
}

// PostDatastream todo
func (gdb *GostDatabase) PostDatastream(d entities.Datastream) (*entities.Datastream, error) {
	return nil, nil
}

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

// GetObservation todo
func (gdb *GostDatabase) GetObservation(id string) (*entities.Observation, error) {
	return nil, nil
}

// GetObservations todo
func (gdb *GostDatabase) GetObservations() ([]*entities.Observation, error) {
	return nil, nil
}

// PostObservation todo
func (gdb *GostDatabase) PostObservation(o entities.Observation) (*entities.Observation, error) {
	return nil, nil
}

// ThingExists checks if a thing is present in the database based on a given id
func (gdb *GostDatabase) ThingExists(thingID int) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM %s.thing WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, thingID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}

// LocationExists checks if a location is present in the database based on a given id
func (gdb *GostDatabase) LocationExists(locationID int) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM  %s.location WHERE id = $1 LIMIT 1)", gdb.Schema)
	err := gdb.Db.QueryRow(sql, locationID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}
