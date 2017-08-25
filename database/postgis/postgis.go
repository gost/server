package postgis

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"

	"encoding/json"
	"strconv"
	"strings"

	"github.com/gost/godata"
	gostErrors "github.com/gost/server/errors"
	gostLog "github.com/gost/server/log"
	"github.com/gost/server/sensorthings/models"
	_ "github.com/lib/pq" // postgres driver
	log "github.com/sirupsen/logrus"
)

const (
	//TimeFormat describes the format in which we want our DateTime to display
	TimeFormat = "YYYY-MM-DD\"T\"HH24:MI:SS.MSZ"
)

var logger *log.Entry

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
	QueryBuilder *QueryBuilder
}

func setupLogger() {
	l, err := gostLog.GetLoggerInstance()
	if err != nil {
		log.Error(err)
	}

	logger = l.WithFields(log.Fields{"package": "gost.server.database.postgis"})
}

// NewDatabase initialises the PostgreSQL database
//	host = TCP host:port or Unix socket depending on Network.
//	user = database user
//	password = database password
//	database = name of database
//	ssl = Whether to use secure TCP/IP connections (TLS).
func NewDatabase(host string, port int, user string, password string, database string, schema string, ssl bool, maxIdeConns int, maxOpenConns int, maxTop int) models.Database {
	setupLogger()
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
		QueryBuilder: CreateQueryBuilder(schema, maxTop),
	}
}

// Start the database
func (gdb *GostDatabase) Start() {
	//gdb.QueryBuilder.Test()

	//ToDo: implement SSL
	logger.Infof("Creating database connection, host: %v, port: %v user: %v, database: %v, schema: %v ssl: %v", gdb.Host, gdb.Port, gdb.User, gdb.Database, gdb.Schema, gdb.Ssl)
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", gdb.Host, gdb.User, gdb.Password, gdb.Database)
	db, err := sql.Open("postgres", dbInfo)
	db.SetMaxIdleConns(gdb.MaxIdeConns)
	db.SetMaxOpenConns(gdb.MaxOpenConns)

	if err != nil {
		logger.Fatal(err)
	}

	gdb.Db = db
	logger.Infof("Connected to database")
}

// CreateSchema creates the needed schema in the database
func (gdb *GostDatabase) CreateSchema(location string) error {
	create, err := GetCreateDatabaseQuery(location, gdb.Schema)
	if err != nil {
		return err
	}

	c := *create
	_, err2 := gdb.Db.Exec(c)
	return err2
}

// GetCreateDatabaseQuery returns the database creation script for PostgreSQL
func GetCreateDatabaseQuery(location string, schema string) (*string, error) {
	bytes, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}

	content := string(bytes[:])
	formatted := fmt.Sprintf(content, schema, schema, schema, schema)
	return &formatted, nil
}

// ContainsToLower checks a string array, array and given string are set to lower-case
func ContainsToLower(s []*godata.SelectItem, e string) bool {
	for _, a := range s {
		for _, b := range a.Segments {
			if strings.ToLower(b.Value) == strings.ToLower(e) {
				return true
			}
		}
	}
	return false
}

// EntityExists checks if entity exists in database
func EntityExists(gdb *GostDatabase, id interface{}, entityName string) bool {
	var result bool
	sql := fmt.Sprintf("SELECT exists (SELECT 1 FROM %s.%s WHERE id = $1 LIMIT 1)", gdb.Schema, entityName)
	err := gdb.Db.QueryRow(sql, id).Scan(&result)
	if err != nil {
		return false
	}

	return result
}

// DeleteEntity deletes a record from database for entity
func DeleteEntity(gdb *GostDatabase, id interface{}, entityName string) error {
	intID, ok := ToIntID(id)
	if !ok {
		errorMessage := fmt.Sprintf("%s does not exist", entityName)
		return gostErrors.NewRequestNotFound(errors.New(errorMessage))
	}

	r, err := gdb.Db.Exec(fmt.Sprintf("DELETE FROM %s.%s WHERE id = $1", gdb.Schema, entityName), intID)
	if err != nil {
		return err
	}

	if c, _ := r.RowsAffected(); c == 0 {
		errorMessage := fmt.Sprintf("%s not found", entityName)
		return gostErrors.NewRequestNotFound(errors.New(errorMessage))
	}
	return nil
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

// ToIntID converts an interface to int id used for the id's in the database
func ToIntID(id interface{}) (int, bool) {
	switch t := id.(type) {
	case string:
		intID, err := strconv.Atoi(t)
		if err != nil {
			return 0, false
		}
		return intID, true
	case float64:
		return int(t), true
	}

	intID, err := strconv.Atoi(fmt.Sprintf("%v", id))
	if err != nil {
		// why not return:  0, err
		return 0, false
	}

	// why not return: intID, nil
	return intID, true
}

func (gdb *GostDatabase) updateEntityColumns(table string, updates map[string]interface{}, entityID int) error {
	if len(updates) == 0 {
		return nil
	}

	columns := ""
	prefix := ""
	for k, v := range updates {
		switch t := v.(type) {
		case string:
			// do not format when value contains ST_ at the start
			if !strings.HasPrefix(t, "ST_") {
				v = fmt.Sprintf("'%s'", t)
			}
		}

		columns += fmt.Sprintf("%s%s=%v", prefix, k, v)
		if prefix != ", " {
			prefix = ", "
		}
	}

	sql := fmt.Sprintf("update %s.%s set %s where id = $1", gdb.Schema, table, columns)
	_, err := gdb.Db.Exec(sql, entityID)
	return err
}
