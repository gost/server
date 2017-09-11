package postgis

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	entities "github.com/gost/core"
	gostLog "github.com/gost/server/log"
	log "github.com/sirupsen/logrus"
)

var idAsSuffix = fmt.Sprintf("%s%s", asSeparator, idField)

// ExecuteSelectCount runs a given count query and returns the value
func ExecuteSelectCount(db *sql.DB, sql string) (int, error) {
	if logger.Logger.Level == log.DebugLevel {
		defer gostLog.DebugfWithElapsedTime(logger, time.Now(), "execute count query: %s", sql)
	}

	var count int
	db.QueryRow(sql).Scan(&count)

	return count, nil
}

// ExecuteSelect executes the select query and creates the retrieved entities
func ExecuteSelect(db *sql.DB, q *QueryParseInfo, sql string) ([]entities.Entity, error) {
	if logger.Logger.Level == log.DebugLevel {
		defer gostLog.DebugfWithElapsedTime(logger, time.Now(), "execute select query: %s", sql)
	}

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valueP := make([]interface{}, count)
	currentQIDEntityID := make(map[int]interface{})                    // keeps track of the current query id and entity id
	parentEntities := []entities.Entity{}                              // array of parent entities
	subEntities := map[int]map[int]map[interface{}][]entities.Entity{} // map of sub entities with a relation to their parent entity map[qid]map[paren qid]map[parent entity id]map[entity id]entity
	relationMap := q.GetQueryIDRelationMap(nil)
	parsedMap := make(map[int]map[int]map[interface{}]map[interface{}]interface{})

	queryParseInfoMap, deleteIDMap, ranges, asMap := prepareExecuteSelect(columns, q)

	for rows.Next() {
		for i := range columns {
			valueP[i] = &values[i]
		}

		rows.Scan(valueP...)

		// split a row into the desired entities (row can contain multiple entities due to join queries)
		sortEntities := rowToEntities(columns, values, ranges, asMap)

		keys := []int{}
		for k := range sortEntities {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		// for every entity found in a row
		for _, k := range keys {
			qi := k
			data := sortEntities[k]
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
						for _, e := range parentEntities {
							if e.GetID() == val {
								skip = true
								break
							}
						}
					} else {
						_, skip = parsedMap[qi][relationMap[qi]][currentQIDEntityID[relationMap[qi]]][val]
						if !skip {
							_, ok := subEntities[qi]
							if !ok {
								subEntities[qi] = make(map[int]map[interface{}][]entities.Entity)
								parsedMap[qi] = make(map[int]map[interface{}]map[interface{}]interface{})
							}
							_, ok = subEntities[qi][relationMap[qi]]
							if !ok {
								subEntities[qi][relationMap[qi]] = make(map[interface{}][]entities.Entity)
								parsedMap[qi][relationMap[qi]] = make(map[interface{}]map[interface{}]interface{})
							}
							_, ok = subEntities[qi][relationMap[qi]][currentQIDEntityID[relationMap[qi]]]
							if !ok {
								subEntities[qi][relationMap[qi]][currentQIDEntityID[relationMap[qi]]] = make([]entities.Entity, 0)
								parsedMap[qi][relationMap[qi]][currentQIDEntityID[relationMap[qi]]] = make(map[interface{}]interface{})
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
				parentEntities = append(parentEntities, newEntity)
			} else {
				subEntities[qi][relationMap[qi]][currentQIDEntityID[relationMap[qi]]] = append(subEntities[qi][relationMap[qi]][currentQIDEntityID[relationMap[qi]]], newEntity)
				parsedMap[qi][relationMap[qi]][currentQIDEntityID[relationMap[qi]]][newEntity.GetID()] = nil
			}
		}
	}

	// Parse all database results into entities starting with parent entities
	for _, entity := range parentEntities {
		parseResults(entity, 0, relationMap, subEntities, deleteIDMap)
	}

	// Remove all needed id fields that weren't requested by the user
	deleteIDFields(parentEntities, deleteIDMap, subEntities)

	// if no parent entities are found return nil, nil
	if len(parentEntities) == 0 {
		return nil, nil
	}

	return parentEntities, nil
}

func prepareExecuteSelect(columns []string, q *QueryParseInfo) (map[int]*QueryParseInfo, map[int]bool, map[int]*QueryParseInfo, map[string]string) {
	queryID := -1
	qpi := q
	queryParseInfoMap := make(map[int]*QueryParseInfo) // QueryID to QueryParseInfo
	deleteIDMap := make(map[int]bool)
	ranges := map[int]*QueryParseInfo{}
	asMap := make(map[string]string)

	// for every _id found store the QueryParseInfo so we know where the column belongs to an create asMap

	for i, c := range columns {
		if strings.HasSuffix(c, idAsSuffix) {
			queryID++
			qpi = q.GetQueryParseInfoByQueryIndex(queryID)
			queryParseInfoMap[queryID] = qpi

			// construct deleteIDMap
			deleteIDMap[qpi.QueryIndex] = false
			if qpi.ExpandItem != nil && qpi.ExpandItem.Select != nil {
				found := false

				for _, p := range qpi.ExpandItem.Select.SelectItems {
					if p.Segments[0].Value == "id" {
						found = true
					}
				}
				deleteIDMap[qpi.QueryIndex] = !found
			}
		}
		ranges[i] = qpi

		slashIndex := strings.Index(c, "_")
		asMap[c] = c[slashIndex+1:]
	}

	return queryParseInfoMap, deleteIDMap, ranges, asMap
}

func rowToEntities(columns []string, values []interface{}, ranges map[int]*QueryParseInfo, asMap map[string]string) map[int]map[string]interface{} {
	sortEntities := make(map[int]map[string]interface{})

	for i := range columns {
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

		sortEntities[ranges[i].QueryIndex][asMap[columns[i]]] = v
	}

	return sortEntities
}

func deleteIDFields(parentEntities []entities.Entity, deleteIDMap map[int]bool, subEntities map[int]map[int]map[interface{}][]entities.Entity) {
	for i := range parentEntities {
		if deleteIDMap[0] {
			parentEntities[i].SetID(nil)
		}

		for subQI := range subEntities {
			if !deleteIDMap[subQI] {
				continue
			}

			for parentQID := range subEntities[subQI] {
				for parentID := range subEntities[subQI][parentQID] {
					entityMap := subEntities[subQI][parentQID][parentID]
					for j := range entityMap {
						entityMap[j].SetID(nil)
					}
				}
			}

			parentIDMap := subEntities[subQI]
			for pID := range parentIDMap {
				for k := range parentIDMap[pID] {
					for l := range parentIDMap[pID][k] {
						parentIDMap[pID][k][l].SetID(nil)
					}
				}
			}
		}
	}
}

func parseResults(entity entities.Entity, from int, relationMap map[int]int, subEntities map[int]map[int]map[interface{}][]entities.Entity, removeIDMap map[int]bool) {
	for subQI := range subEntities {
		relation, ok := relationMap[subQI]
		if ok && relation == from {
			relatedEntityMap, ok := subEntities[subQI][from][entity.GetID()]
			if ok {
				addRelationToEntity(entity, relatedEntityMap)
				for _, relatedEntity := range relatedEntityMap {
					parseResults(relatedEntity, subQI, relationMap, subEntities, removeIDMap)
				}
			}
		}
	}
}

func addRelationToEntity(parent entities.Entity, subEntities []entities.Entity) {
	switch parentEntity := parent.(type) {
	case *entities.Thing:
		addRelationToThing(parentEntity, subEntities)
	case *entities.Location:
		addRelationToLocation(parentEntity, subEntities)
	case *entities.HistoricalLocation:
		addRelationToHistoricalLocation(parentEntity, subEntities)
	case *entities.Datastream:
		addRelationToDatastream(parentEntity, subEntities)
	case *entities.Sensor:
		addRelationToSensor(parentEntity, subEntities)
	case *entities.ObservedProperty:
		addRelationToObservedProperty(parentEntity, subEntities)
	case *entities.Observation:
		addRelationToObservation(parentEntity, subEntities)
	case *entities.FeatureOfInterest:
		addRelationToFeatureOfInterest(parentEntity, subEntities)
	}
}

func addRelationToThing(parentEntity *entities.Thing, subEntities []entities.Entity) {
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
}

func addRelationToLocation(parentEntity *entities.Location, subEntities []entities.Entity) {
	for _, se := range subEntities {
		switch subEntity := se.(type) {
		case *entities.HistoricalLocation:
			parentEntity.HistoricalLocations = append(parentEntity.HistoricalLocations, subEntity)
		case *entities.Thing:
			parentEntity.Things = append(parentEntity.Things, subEntity)
		}
	}
}

func addRelationToHistoricalLocation(parentEntity *entities.HistoricalLocation, subEntities []entities.Entity) {
	for _, se := range subEntities {
		switch subEntity := se.(type) {
		case *entities.Thing:
			parentEntity.Thing = subEntity
		case *entities.Location:
			parentEntity.Locations = append(parentEntity.Locations, subEntity)
		}
	}
}

func addRelationToDatastream(parentEntity *entities.Datastream, subEntities []entities.Entity) {
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
}

func addRelationToSensor(parentEntity *entities.Sensor, subEntities []entities.Entity) {
	for _, se := range subEntities {
		switch subEntity := se.(type) {
		case *entities.Datastream:
			parentEntity.Datastreams = append(parentEntity.Datastreams, subEntity)
		}
	}
}

func addRelationToObservedProperty(parentEntity *entities.ObservedProperty, subEntities []entities.Entity) {
	for _, se := range subEntities {
		switch subEntity := se.(type) {
		case *entities.Datastream:
			parentEntity.Datastreams = append(parentEntity.Datastreams, subEntity)
		}
	}
}

func addRelationToObservation(parentEntity *entities.Observation, subEntities []entities.Entity) {
	for _, se := range subEntities {
		switch subEntity := se.(type) {
		case *entities.Datastream:
			parentEntity.Datastream = subEntity
		case *entities.FeatureOfInterest:
			parentEntity.FeatureOfInterest = subEntity
		}
	}
}

func addRelationToFeatureOfInterest(parentEntity *entities.FeatureOfInterest, subEntities []entities.Entity) {
	for _, se := range subEntities {
		switch subEntity := se.(type) {
		case *entities.Observation:
			parentEntity.Observations = append(parentEntity.Observations, subEntity)
		}
	}
}
