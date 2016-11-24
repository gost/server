package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
	"strings"
)

// tables as defined in postgis
var (
	thingTable              = "thing"
	locationTable           = "location"
	historicalLocationTable = "historicallocation"
	sensorTable             = "sensor"
	observedPropertyTable   = "observedproperty"
	datastreamTable         = "datastream"
	observationTable        = "observation"
	featureOfInterestTable  = "featureofinterest"
)

// thing fields
var (
	thingID          = fmt.Sprintf("%s.id", thingTable)
	thingName        = fmt.Sprintf("%s.name", thingTable)
	thingDescription = fmt.Sprintf("%s.description", thingTable)
	thingProperties  = fmt.Sprintf("%s.properties", thingTable)
)

// datastream fields
var (
	datastreamID                 = fmt.Sprintf("%s.id", datastreamTable)
	datastreamName               = fmt.Sprintf("%s.name", datastreamTable)
	datastreamDescription        = fmt.Sprintf("%s.description", datastreamTable)
	datastreamUnitOfMeasurement  = fmt.Sprintf("%s.unitofmeasurement", datastreamTable)
	datastreamObservationType    = fmt.Sprintf("%s.observationtype", datastreamTable)
	datastreamObservedArea       = fmt.Sprintf("%s.observedarea", datastreamTable)
	datastreamPhenomenonTime     = fmt.Sprintf("%s.phenomenontime", datastreamTable)
	datastreamResultTime         = fmt.Sprintf("%s.resulttime", datastreamTable)
	datastreamThingID            = fmt.Sprintf("%s.thing_id", datastreamTable)
	datastreamSensorID           = fmt.Sprintf("%s.sensor_id", datastreamTable)
	datastreamObservedPropertyID = fmt.Sprintf("%s.observedproperty_id", datastreamTable)
)

// maps an entity property name to the right field
var selectMapping = map[entities.EntityType]map[string]string{
	entities.EntityTypeThing: {
		"id":          thingID,
		"name":        thingName,
		"description": thingDescription,
		"properties":  thingProperties,
	},
	entities.EntityTypeDatastream: {
		"id":                  datastreamID,
		"name":                datastreamName,
		"description":         datastreamDescription,
		"unitofmeasurement":   datastreamUnitOfMeasurement,
		"observationtype":     datastreamObservationType,
		"observedarea":        fmt.Sprintf("public.ST_AsGeoJSON(%s)", datastreamObservedArea),
		"phenomenontime":      datastreamPhenomenonTime,
		"resulttime":          datastreamResultTime,
		"thing_id":            datastreamThingID,
		"sensor_id":           datastreamSensorID,
		"observedproperty_id": datastreamObservedPropertyID,
	},
}

type queryBuilder struct {
	maxTop int
	schema string
	tables map[entities.EntityType]string
	joins  map[entities.EntityType]map[entities.EntityType]string
}

// CreateQueryBuilder instantiates a new queryBuilder, the queryBuilder is used to create
// select queries based on the given entities, id en QueryOptions (ODATA)
// schema is the used database schema can be empty, maxTop is the maximum top the query should return
func CreateQueryBuilder(schema string, maxTop int) *queryBuilder {
	qb := &queryBuilder{
		schema: schema,
		maxTop: maxTop,
		tables: createEntityTableMap(schema),
		joins:  createJoinMapping(),
	}

	return qb
}

func createEntityTableMap(schema string) map[entities.EntityType]string {
	if len(schema) > 0 && !strings.Contains(schema, ".") {
		schema = fmt.Sprintf("%s.", schema)
	}

	entityTypeTableMap := map[entities.EntityType]string{
		entities.EntityTypeThing:              fmt.Sprintf("%s%s", schema, thingTable),
		entities.EntityTypeLocation:           fmt.Sprintf("%s%s", schema, locationTable),
		entities.EntityTypeHistoricalLocation: fmt.Sprintf("%s%s", schema, historicalLocationTable),
		entities.EntityTypeSensor:             fmt.Sprintf("%s%s", schema, sensorTable),
		entities.EntityTypeObservedProperty:   fmt.Sprintf("%s%s", schema, observedPropertyTable),
		entities.EntityTypeDatastream:         fmt.Sprintf("%s%s", schema, datastreamTable),
		entities.EntityTypeObservation:        fmt.Sprintf("%s%s", schema, observationTable),
		entities.EntityTypeFeatureOfInterest:  fmt.Sprintf("%s%s", schema, featureOfInterestTable),
	}

	return entityTypeTableMap
}

func createJoinMapping() map[entities.EntityType]map[entities.EntityType]string {
	joinMap := map[entities.EntityType]map[entities.EntityType]string{
		entities.EntityTypeDatastream: {
			entities.EntityTypeThing: fmt.Sprintf("%s = %s", datastreamThingID, thingID),
		},
		entities.EntityTypeThing: {
			entities.EntityTypeDatastream: fmt.Sprintf("%s = %s", thingID, datastreamThingID),
		},
	}

	return joinMap
}

// toAs returns the AS name for a field
func (qb *queryBuilder) toAs(field string) string {
	return strings.Replace(field, ".", "_", -1)
}

// removeSchema removes the prefix in front of a table
func (qb *queryBuilder) removeSchema(table string) string {
	i := strings.Index(table, ".")
	if i == -1 {
		return table
	}

	return table[i+1:]
}

// getLimit returns the max entities to retrieve, this number is set by ODATA's
// $top, if not provided use the global value
func (qb *queryBuilder) getLimit(qo *odata.QueryOptions) string {
	if qo != nil && !qo.QueryTop.IsNil() {
		return fmt.Sprintf("%v", qo.QueryTop.Limit)
	} else {
		return fmt.Sprintf("%v", qb.maxTop)
	}
}

// getOffset returns the offset, this number is set by ODATA's
// $skip, if not provided do not skip anything = return "0"
func (qb *queryBuilder) getOffset(qo *odata.QueryOptions) string {
	if qo != nil && !qo.QueryTop.IsNil() {
		return fmt.Sprintf("%v", qo.QueryTop.Limit)
	} else {
		return "0"
	}
}

// getOrderBy returns the string that needs to be placed after ORDER BY, this is set using
// ODATA's $orderby if not given use the default ORDER BY "table".id DESC
func (qb *queryBuilder) getOrderBy(et entities.EntityType, qo *odata.QueryOptions) string {
	if qo != nil && !qo.QueryOrderBy.IsNil() {
		return fmt.Sprintf("%v %v", selectMapping[et][qo.QueryOrderBy.Property], strings.ToUpper(qo.QueryOrderBy.Suffix))
	}

	return fmt.Sprintf("%s DESC", selectMapping[et]["id"])
}

// ToDo: implement filter
// getFilter returns a string that can be placed after WHERE, the filter is set by
// ODATA's $filter and can be found in qo.QueryFilter
func (qb *queryBuilder) getFilter(et entities.EntityType, qo *odata.QueryOptions) string {
	if qo != nil && !qo.QueryFilter.IsNil() {
		return fmt.Sprintf("%v", "")
	}

	return ""
}

// getSelect return the select string that needs to be placed after SELECT in the query
// select is set by ODATA's $select, if not set get all properties for the given entity (return all)
// addID to true if it needs to be added and isn't in QuerySelect.Params, addAs to true if a field needs to be
// outputted with AS [name]
func (qb *queryBuilder) getSelect(et entities.Entity, qo *odata.QueryOptions, addID bool, addAs bool, selectString string) string {
	var properties []string
	if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
		properties = et.GetPropertyNames()
	} else {
		idAdded := false
		for _, p := range qo.QuerySelect.Params {
			if p == "id" {
				idAdded = true
			}
			for _, pn := range et.GetPropertyNames() {
				if strings.ToLower(p) == strings.ToLower(pn) {
					properties = append(properties, pn)
				}
			}
		}
		if addID && !idAdded {
			properties = append(properties, "id")
		}
	}

	for _, p := range properties {
		toAdd := ""
		if len(selectString) > 0 {
			toAdd += ", "
		}

		field := selectMapping[et.GetEntityType()][p]
		if addAs {
			selectString += fmt.Sprintf("%s%s as %s", toAdd, field, qb.toAs(field))
		} else {
			selectString += fmt.Sprintf("%s%s", toAdd, field)
		}

	}

	if qo != nil && !qo.QueryExpand.IsNil() {
		for _, o := range qo.QueryExpand.Operations {
			selectString = qb.getSelect(o.Entity, o.QueryOptions, addID, addAs, selectString)
		}
	}

	return selectString
}

func (qb *queryBuilder) createLateralJoin(e1 entities.Entity, e2 entities.Entity, qo *odata.QueryOptions, joinString string) string {
	if e2 != nil {
		et2 := e2.GetEntityType()
		joinString = fmt.Sprintf("%s "+
			"INNER JOIN LATERAL ("+
			"SELECT %s FROM %s WHERE %s "+
			"ORDER BY %s "+
			"LIMIT %s OFFSET %s) as %s on true", joinString,
			qb.getSelect(e2, qo, true, false, ""),
			qb.tables[et2],
			qb.joins[e1.GetEntityType()][et2],
			qb.getOrderBy(et2, qo),
			qb.getLimit(qo),
			qb.getOffset(qo),
			qb.removeSchema(qb.tables[et2]))
	} else {
		if qo != nil && !qo.QueryExpand.IsNil() {
			for _, qe := range qo.QueryExpand.Operations {
				joinString = qb.createLateralJoin(e1, qe.Entity, qe.QueryOptions, joinString)
			}
		}
	}

	return joinString
}

// CreateQuery creates a new query based on given input
//   e1 = entity to get
//   e2 = from entity
//   id = where e2
// example: Datastreams(1)/Thing = CreateQuery(&entities.Thing, &entities.Datastream, 1, nil)
func (qb *queryBuilder) CreateQuery(e1 entities.Entity, e2 entities.Entity, id interface{}, qo *odata.QueryOptions) (string, error) {
	et1 := e1.GetEntityType()
	et2 := e1.GetEntityType()
	if e2 != nil {
		et2 = e2.GetEntityType()
	}

	queryString := fmt.Sprintf("SELECT %s FROM %s %s", qb.getSelect(e1, qo, false, true, ""), qb.tables[et1], qb.createLateralJoin(e1, e2, qo, ""))
	if id != nil {
		queryString = fmt.Sprintf("%s WHERE %s = %v", queryString, selectMapping[et2]["id"], id)
	}

	queryString = fmt.Sprintf("%s ORDER BY %s", queryString, qb.getOrderBy(et1, qo))
	queryString = fmt.Sprintf("%s LIMIT %s OFFSET %s", queryString, qb.getLimit(qo), qb.getOffset(qo))
	return queryString, nil
}

func (qb *queryBuilder) Test() {
	sql1, _ := qb.CreateQuery(&entities.Thing{}, nil, nil, nil)
	fmt.Println(sql1)
	fmt.Println("------------------------")

	qo2 := &odata.QueryOptions{}
	qo2.QuerySelect = &odata.QuerySelect{}
	qo2.QuerySelect.Parse("name,description")
	sql2, _ := qb.CreateQuery(&entities.Thing{}, &entities.Datastream{}, 1, qo2)
	fmt.Println(sql2)
	fmt.Println("------------------------")

	qo31 := &odata.QueryOptions{}
	qo31.QuerySelect = &odata.QuerySelect{}
	qo31.QuerySelect.Parse("name,description")
	qo31.QueryExpand = &odata.QueryExpand{}
	qo31.QueryExpand.Parse("Thing($select=name)")
	sql3, _ := qb.CreateQuery(&entities.Datastream{}, nil, nil, qo31)
	fmt.Println(sql3)
	fmt.Println("------------------------")
}
