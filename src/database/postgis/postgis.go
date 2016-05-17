package postgis

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"encoding/json"
	"github.com/geodan/gost/src/sensorthings/models"
	_ "github.com/lib/pq" // postgres driver
	"strings"
)

const (
	//TimeFormat describes the format in which we want our DateTime to display
	TimeFormat = "YYYY-MM-DD\"T\"HH24:MI:SS.MSZ"
)

// GostDatabase implementation
type GostDatabase struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
	Schema       string
	Ssl          bool
	MaxIdeConns  int
	MaxOpenConns int
	Db           *sql.DB
}

// NewDatabase initialises the PostgreSQL database
//	host = TCP host:port or Unix socket depending on Network.
//	user = database user
//	password = database password
//	database = name of database
//	ssl = Whether to use secure TCP/IP connections (TLS).
func NewDatabase(host string, port int, user string, password string, database string, schema string, ssl bool, maxIdeConns int, maxOpenConns int) models.Database {
	return &GostDatabase{
		Host:         host,
		Port:         port,
		User:         user,
		Password:     password,
		Database:     database,
		Schema:       schema,
		Ssl:          ssl,
		MaxIdeConns:  maxIdeConns,
		MaxOpenConns: maxOpenConns,
	}
}

// Start the database
func (gdb *GostDatabase) Start() {
	//ToDo: implement SSL
	log.Println("Creating database connection...")
	log.Printf("Database, host: \"%v\", port: \"%v\" user: \"%v\", database: \"%v\", schema: \"%v\" ssl: \"%v\"", gdb.Host, gdb.Port, gdb.User, gdb.Database, gdb.Schema, gdb.Ssl)

	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", gdb.Host, gdb.User, gdb.Password, gdb.Database)
	db, err := sql.Open("postgres", dbInfo)
	db.SetMaxIdleConns(gdb.MaxIdeConns)
	db.SetMaxOpenConns(gdb.MaxOpenConns)

	if err != nil {
		log.Fatal(err)
	}

	gdb.Db = db
	err2 := gdb.Db.Ping()
	if err2 != nil {
		log.Fatal("Unable to connect to database, check your network connection.")
	}

	log.Printf("Connected to database, host: \"%v\", port: \"%v\" user: \"%v\", database: \"%v\", schema: \"%v\" ssl: \"%v\"", gdb.Host, gdb.Port, gdb.User, gdb.Database, gdb.Schema, gdb.Ssl)
}

// CreateSchema creates the needed schema in the database
func (gdb *GostDatabase) CreateSchema(location string) error {
	create, err := GetCreateDatabaseQuery(location, gdb.Schema)
	if err != nil {
		return err
	}

	c := *create
	_, err2 := gdb.Db.Exec(c)
	if err2 != nil {
		return err2
	}

	return nil
}

// GetCreateDatabaseQuery returns the database creation script for PostgreSQL
func GetCreateDatabaseQuery(location string, schema string) (*string, error) {
	bytes, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}

	content := string(bytes[:])
	formatted := fmt.Sprintf(content, schema, schema)
	return &formatted, nil
}

// JSONToMap converts a string of json into a map
func JSONToMap(data *string) (map[string]interface{}, error) {
	var p map[string]interface{}
	if data == nil || len(*data) == 0 {
		return p, nil
	}

	err := json.Unmarshal([]byte(*data), &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// ConvertNullString converts a scanned value from ro into a usable string
func ConvertNullString(value *string) string {
	if value == nil {
		return ""
	}

	return fmt.Sprintf("%s", *value)
}

// PrepareTimeRangeForPostgres splits an incoming timerange by the / delimiter and returns the start
// and end value, if it can't be splitted the function will return the start value also as end
func PrepareTimeRangeForPostgres(timeRange string) string {
	if len(timeRange) == 0 {
		return "NULL"
	}

	s := strings.Split(timeRange, "/")
	if len(s) == 1 {
		return fmt.Sprintf("tstzrange('%s', NULL)", s[0])
	}

	return fmt.Sprintf("tstzrange('%s','%s')", s[0], s[1])
}

// TimeRangeToString converts a start and end datetime string into a string valid time string (SensorThings spec)
func TimeRangeToString(start, end *string) string {
	if start != nil && end == nil {
		return fmt.Sprintf("%s", *start)
	}
	if start != nil && end != nil {
		return fmt.Sprintf("%s/%s", *start, *end)
	}

	return ""
}
