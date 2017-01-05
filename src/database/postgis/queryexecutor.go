package postgis

import (
	"database/sql"
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"

	"strings"
)

var idAsSuffix = fmt.Sprintf("%s%s", asSeparator, idField)

// ExecuteSelectCount runs a given count query and returns the value
func ExecuteSelectCount(db *sql.DB, sql string) (int, error) {
	var count int
	db.QueryRow(sql).Scan(&count)

	return count, nil
}

// ExecuteSelect executes the select query and creates the retrieved entities
func ExecuteSelect(db *sql.DB, q *QueryParseInfo, sql string) ([]interface{}, error) {
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valueP := make([]interface{}, count)
	queryParseInfoMap := make(map[int]*QueryParseInfo)
	currentQIDEntityID := make(map[int]interface{})                                  // keeps track of the current query id and entity id
	parentEntities := map[interface{}]entities.Entity{}                              // array of parent entities
	subEntities := map[int]map[int]map[interface{}]map[interface{}]entities.Entity{} // map of sub entities with a relation to their parent entity map[qid]map[paren qid]map[parent entity id]map[entity id]entity
	relationMap := q.GetQueryIDRelationMap(nil)

	// for every _id found store the QueryParseInfo so we know where the column belongs to
	ranges := map[int]*QueryParseInfo{}
	qpi := q
	queryId := -1
	for i, c := range columns {
		if strings.HasSuffix(c, idAsSuffix) {
			queryId++
			qpi = q.GetQueryParseInfoByQueryIndex(queryId)
			queryParseInfoMap[queryId] = qpi
		}
		ranges[i] = qpi
	}

	for rows.Next() {
		for i := range columns {
			valueP[i] = &values[i]
		}

		rows.Scan(valueP...)
		sortEntities := make(map[int]map[string]interface{})

		// split a row into the desired entities (row can contain multiple entities due to join queries)
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

		// for every entity found in a row
		for qi, data := range sortEntities {
			// filter out already parsed and nil entities
			skip := false
			for col, val := range data {
				if strings.HasSuffix(col, idAsSuffix) {
					if val == nil {
						skip = true
						break
					} else {
						currentQIDEntityID[qi] = val
					}

					if qi == 0 {
						_, skip = parentEntities[val]
					} else {
						_, skip = subEntities[qi][relationMap[qi]][currentQIDEntityID[relationMap[qi]]][val]
						if !skip {
							_, ok := subEntities[qi]
							if !ok {
								subEntities[qi] = make(map[int]map[interface{}]map[interface{}]entities.Entity)
							}
							_, ok = subEntities[qi][relationMap[qi]]
							if !ok {
								subEntities[qi][relationMap[qi]] = make(map[interface{}]map[interface{}]entities.Entity)
							}
							_, ok = subEntities[qi][relationMap[qi]][currentQIDEntityID[relationMap[qi]]]
							if !ok {
								subEntities[qi][relationMap[qi]][currentQIDEntityID[relationMap[qi]]] = make(map[interface{}]entities.Entity)
							}
						}
					}

					break
				}
			}

			// entity already parsed continue with next one
			if skip {
				continue
			}

			newEntity, err := queryParseInfoMap[qi].Parse(data)
			if err != nil {
				return nil, err
			}

			if qi == 0 {
				parentEntities[currentQIDEntityID[qi]] = newEntity
			} else {
				subEntities[qi][relationMap[qi]][currentQIDEntityID[relationMap[qi]]][currentQIDEntityID[qi]] = newEntity
			}
		}
	}

	for _, entity := range parentEntities {
		parseResults(entity, 0, relationMap, subEntities)
	}

	parentEntitiesLength := len(parentEntities)
	if parentEntitiesLength == 0 {
		return nil, nil
	} else {
		entitySlice := make([]interface{}, 0)
		for _, e := range parentEntities {
			entitySlice = append(entitySlice, e)
		}

		return entitySlice, nil
	}

	return nil, nil
}

func parseResults(entity entities.Entity, from int, relationMap map[int]int, subEntities map[int]map[int]map[interface{}]map[interface{}]entities.Entity) {
	for subQI := range subEntities {
		relation, ok := relationMap[subQI]
		if ok && relation == from {
			relatedEntityMap, ok := subEntities[subQI][from][entity.GetID()]
			if ok {
				addRelationToEntity(entity, relatedEntityMap)
				for _, relatedEntity := range relatedEntityMap {
					parseResults(relatedEntity, subQI, relationMap, subEntities)
				}
			}
		}
	}
}

func addRelationToEntity(parent entities.Entity, subEntities map[interface{}]entities.Entity) {
	switch parentEntity := parent.(type) {
	case *entities.Thing:
		for _, se := range subEntities {
			switch subEntity := se.(type) {
			case *entities.HistoricalLocation:
				parentEntity.HistoricalLocations = append(parentEntity.HistoricalLocations, subEntity)
			case *entities.Location:
				parentEntity.Locations = append(parentEntity.Locations, subEntity)
			case *entities.Datastream:
				parentEntity.Datastreams = append(parentEntity.Datastreams, subEntity)
			}
		}
	case *entities.Location:
		for _, se := range subEntities {
			switch subEntity := se.(type) {
			case *entities.HistoricalLocation:
				parentEntity.HistoricalLocations = append(parentEntity.HistoricalLocations, subEntity)
			case *entities.Thing:
				parentEntity.Things = append(parentEntity.Things, subEntity)
			}
		}
	case *entities.HistoricalLocation:
		for _, se := range subEntities {
			switch subEntity := se.(type) {
			case *entities.Thing:
				parentEntity.Thing = subEntity
			case *entities.Location:
				parentEntity.Locations = append(parentEntity.Locations, subEntity)
			}
		}
	case *entities.Datastream:
		for _, se := range subEntities {
			switch subEntity := se.(type) {
			case *entities.Observation:
				parentEntity.Observations = append(parentEntity.Observations, subEntity)
			case *entities.Thing:
				parentEntity.Thing = subEntity
			case *entities.Sensor:
				parentEntity.Sensor = subEntity
			case *entities.ObservedProperty:
				parentEntity.ObservedProperty = subEntity
			}
		}
	case *entities.Sensor:
		for _, se := range subEntities {
			switch subEntity := se.(type) {
			case *entities.Datastream:
				parentEntity.Datastreams = append(parentEntity.Datastreams, subEntity)
			}
		}
	case *entities.ObservedProperty:
		for _, se := range subEntities {
			switch subEntity := se.(type) {
			case *entities.Datastream:
				parentEntity.Datastreams = append(parentEntity.Datastreams, subEntity)
			}
		}
	case *entities.Observation:
		for _, se := range subEntities {
			switch subEntity := se.(type) {
			case *entities.Datastream:
				parentEntity.Datastream = subEntity
			case *entities.FeatureOfInterest:
				parentEntity.FeatureOfInterest = subEntity
			}
		}
	case *entities.FeatureOfInterest:
		for _, se := range subEntities {
			switch subEntity := se.(type) {
			case *entities.Observation:
				parentEntity.Observations = append(parentEntity.Observations, subEntity)
			}
		}
	}
}
