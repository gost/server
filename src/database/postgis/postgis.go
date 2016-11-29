package postgis

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"encoding/json"
	"strconv"
	"strings"

	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
	_ "github.com/lib/pq" // postgres driver
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
	QueryBuilder *queryBuilder
}

// NewDatabase initialises the PostgreSQL database
//	host = TCP host:port or Unix socket depending on Network.
//	user = database user
//	password = database password
//	database = name of database
//	ssl = Whether to use secure TCP/IP connections (TLS).
func NewDatabase(host string, port int, user string, password string, database string, schema string, ssl bool, maxIdeConns int, maxOpenConns int, maxTop int) models.Database {
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
	log.Printf("Connected to database, host: \"%v\", port: \"%v\" user: \"%v\", database: \"%v\", schema: \"%v\" ssl: \"%v\"", gdb.Host, gdb.Port, gdb.User, gdb.Database, gdb.Schema, gdb.Ssl)

	//TEST
	sql9, _ := gdb.QueryBuilder.CreateQuery(&entities.Observation{}, &entities.Datastream{}, 1, nil)
	ExecuteSelect(gdb.Db, sql9)
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
	formatted := fmt.Sprintf(content, schema, schema, schema, schema)
	return &formatted, nil
}

// Contains checks a string array
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Contains checks a string array, array and given string are set to lower-case
func ContainsToLower(s []string, e string) bool {
	for _, a := range s {
		if strings.ToLower(a) == strings.ToLower(e) {
			return true
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
		return 0, false
	}
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
	if err != nil {
		return err
	}

	return nil
}

// CreateSelectString creates a select string based on available parameters and or QuerySelect option
func CreateSelectString(e entities.Entity, qo *odata.QueryOptions, prefix string, trail string, mapping map[string]string) string {
	var properties []string

	if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
		properties = e.GetPropertyNames()
	} else {
		for _, p := range qo.QuerySelect.Params {
			for _, pn := range e.GetPropertyNames() {
				if strings.ToLower(p) == strings.ToLower(pn) {
					properties = append(properties, pn)
				}
			}
		}
	}

	s := ""
	for _, p := range properties {
		skip := false
		for _, e := range entities.EntityTypeList {
			if p == e.ToString() {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		up := p

		if len(prefix) > 0 {
			p = prefix + p
		}

		if len(trail) > 0 {
			p += trail
		}

		if mapping != nil || len(mapping) > 0 {
			m, ok := mapping[up]
			if ok {
				p = m
			}
		}

		toAdd := ""

		if len(s) > 0 {
			toAdd += ", "
		}

		s += toAdd + p
	}

	return s
}

// CreateTopSkipQueryString creates a LIMIT and OFFSET query string
func CreateTopSkipQueryString(qo *odata.QueryOptions) string {
	q := ""
	if qo != nil && !qo.QueryTop.IsNil() {
		q += fmt.Sprintf(" LIMIT %v", qo.QueryTop.Limit)
	}
	if qo != nil && !qo.QuerySkip.IsNil() {
		q += fmt.Sprintf(" OFFSET %v", qo.QuerySkip.Index)
	}

	return q
}

// CreateFilterQueryString converts an OData query string found in odata.QueryOptions.QueryFilter to a PostgreSQL query string
// ParamFactory is used for converting SensorThings parameter names to postgres field names
// Convert receives a name such as phenomenonTime and returns "data ->> 'id'" true, returns
// false if parameter cannot be converted
func CreateFilterQueryString(qo *odata.QueryOptions, paramFactory func(string, interface{}) (string, string, error), prefix string) (string, error) {
	q := ""
	if qo != nil && !qo.QueryFilter.IsNil() {
		q += prefix
		ps, ops := qo.QueryFilter.Predicate.Split()
		for i, p := range ps {
			var left, right, operator string
			var err error

			if left, right, err = paramFactory(fmt.Sprintf("%v", p.Left), p.Right); err != nil {
				return "", err
			}

			if operator, err = OdataOperatorToPostgreSQL(p.Operator); err != nil {
				return "", err
			}

			q += fmt.Sprintf("%v %v %v", left, operator, right)
			if len(ops)-1 >= i {
				q += fmt.Sprintf(" %v ", ops[i])
			}
		}
		q += " "
	}

	return q, nil
}

// OdataOperatorToPostgreSQL converts an odata.OdataOperator to a PostgreSQL string representation
func OdataOperatorToPostgreSQL(o odata.Operator) (string, error) {
	switch o {
	case odata.And:
		return "AND", nil
	case odata.Or:
		return "OR", nil
	case odata.Not:
		return "NOT", nil
	case odata.Equals:
		return "=", nil
	case odata.NotEquals:
		return "!=", nil
	case odata.GreaterThan:
		return ">", nil
	case odata.GreaterThanOrEquals:
		return ">=", nil
	case odata.LessThan:
		return "<", nil
	case odata.LessThanOrEquals:
		return "<=", nil
	case odata.IsNull:
		return "IS NULL", nil
	}

	return "", fmt.Errorf("Operator %v not implemented", o.ToString())
}
