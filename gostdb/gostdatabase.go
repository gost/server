package gostdb

import (
	"fmt"
	"log"

	"github.com/tebben/gost/sensorthings"
	_ "github.com/lib/pq"
	"database/sql"
	"time"
	"strconv"
	"errors"
	"encoding/json"
)

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

// Initialize the PostgreSQL database
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

func (gdb *GostDatabase) CreateSchema() {
	create := GetCreateDatabaseQuery(gdb.Schema)
	_, err := gdb.Db.Exec(create)
	if err != nil {
		log.Fatal(err)
	}
}

func (gdb *GostDatabase) GetThing(id string) (*sensorthings.Thing, error) {
	intId, err := strconv.Atoi(id)
	if err != nil{
		return nil, err
	}

	var thingId int
	var description string
	var properties string
	err2 := gdb.Db.QueryRow("SELECT * FROM v1.thing WHERE id = $1 LIMIT 1", intId).Scan(&thingId, &description, &properties)

	if err2 != nil {
		return nil, err
	}

	thing := sensorthings.Thing{}
	thing.ID = strconv.Itoa(thingId)
	thing.Description = description

	var p map[string]string
	err3 := json.Unmarshal([]byte(properties), &p)
	if err3 != nil {
		return nil, err3
	}

	thing.Properties = p

	return &thing, nil
}

func (gdb *GostDatabase) GetThings() ([]*sensorthings.Thing, error) {
	rows, err := gdb.Db.Query("SELECT * FROM v1.thing")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	things := make([]*sensorthings.Thing, 0)

	for rows.Next() {
		thing := sensorthings.Thing{}

		var id int
		var description string
		var properties string
		err2 := rows.Scan(&id, &description, &properties)
		if err2 != nil{
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

func (gdb *GostDatabase) PostThing(thing sensorthings.Thing) (*sensorthings.Thing, error) {
	jsonProperties, _ := json.Marshal(thing.Properties)
	var thingId int
	err := gdb.Db.QueryRow("INSERT INTO v1.thing (description, properties) VALUES ($1, $2) RETURNING id", thing.Description, jsonProperties).Scan(&thingId)
	if err != nil {
		return nil, err
	}

	thing.ID = strconv.Itoa(thingId)
	return &thing, nil
}

func (gdb *GostDatabase) PostLocation(location sensorthings.Location) (*sensorthings.Location, error) {
	var locationId int
	err := gdb.Db.QueryRow("INSERT INTO v1.location (description, encodingtype, location) VALUES ($1, $2, $3) RETURNING id", location.Description, 1, location.Location).Scan(&locationId)
	if err != nil {
		return nil, err
	}

	location.ID = strconv.Itoa(locationId)
	return &location, nil
}

func (gdb *GostDatabase) PostHistoricalLocation(thingID string, locationID string) error {
	tid, err := strconv.Atoi(thingID)
	if(!gdb.ThingExists(tid) || err != nil){
		return errors.New(fmt.Sprintf("Thing(%v) does not exist", thingID))
	}

	lid, err2 := strconv.Atoi(thingID)
	if(!gdb.ThingExists(lid) || err2 != nil){
		return errors.New(fmt.Sprintf("Location(%v) does not exist", locationID))
	}

	//check if thing and location exist
	_, err3 := gdb.Db.Exec("INSERT INTO v1.historicallocation (time, thing_id, location_id) VALUES ($1, $2, $3)", time.Now(), tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}


func (gdb *GostDatabase) LinkLocation(thingID string, locationID string) error {
	tid, err := strconv.Atoi(thingID)
	if(!gdb.ThingExists(tid) || err != nil){
		return errors.New(fmt.Sprintf("Thing(%v) does not exist", thingID))
	}

	lid, err2 := strconv.Atoi(thingID)
	if(!gdb.ThingExists(lid) || err2 != nil){
		return errors.New(fmt.Sprintf("Location(%v) does not exist", locationID))
	}

	_, err3 := gdb.Db.Exec("INSERT INTO v1.thing_to_location (thing_id, location_id) VALUES ($1, $2)", tid, lid)
	if err3 != nil {
		return err3
	}

	return nil
}

func (gdb *GostDatabase) ThingExists(thingID int) bool{
	var result bool
	err := gdb.Db.QueryRow("SELECT exists (SELECT 1 FROM v1.thing WHERE id = $1 LIMIT 1)", thingID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}

func (gdb *GostDatabase) LocationExists(locationID int) bool{
	var result bool
	err := gdb.Db.QueryRow("SELECT exists (SELECT 1 FROM v1.location WHERE id = $1 LIMIT 1)", locationID).Scan(&result)
	if err != nil {
		return false
	}

	return result
}