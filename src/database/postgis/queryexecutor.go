package postgis

import (
	"database/sql"
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"strings"
)

// ExecuteSelect executes the select query and creates the retrieved entities
func ExecuteSelect(db *sql.DB, q *QueryParseInfo, sql string) (interface{}, error) {
	fmt.Println(sql)
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valueP := make([]interface{}, count)
	parsed := map[int]map[string]entities.Entity{}

	// for every _id found store the QueryParseInfo so we know where the column belongs to
	ranges := map[int]*QueryParseInfo{}
	qpi := q
	queryId := -1
	for i, c := range columns {
		if strings.HasSuffix(c, "_id") {
			queryId++
			qpi = q.GetQueryParseInfoByQueryIndex(queryId)
		}
		ranges[i] = qpi
	}

	for rows.Next() {
		for i := range columns {
			valueP[i] = &values[i]
		}

		rows.Scan(valueP...)
		sortEntities := make(map[int]map[string]interface{})

		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}

			_, ok = sortEntities[ranges[i].QueryIndex]
			if !ok {
				sortEntities[ranges[i].QueryIndex] = make(map[string]interface{})
			}

			sortEntities[ranges[i].QueryIndex][col] = v
		}

		for k, v := range sortEntities {
			newEntity, err := ranges[k].Parse(v)
			if err != nil {
				return nil, err
			}

			_, ok := parsed[k]
			if !ok {
				parsed[k] = make(map[string]entities.Entity)
			}

			parsed[k][newEntity.GetID().(string)] = newEntity
		}
	}

	// combine expand entities into the correct entity
	// fmt.Println(fmt.Sprintf("%v", newEntity))

	//for

	return nil, nil
}
