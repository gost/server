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
	schema string
	tables map[entities.EntityType]string
	joins  map[entities.EntityType]map[entities.EntityType]string
}

func CreateQueryBuilder(schema string) *queryBuilder {
	qb := &queryBuilder{
		schema: schema,
		tables: createEntityTableMap(schema),
	}

	qb.joins = createJoinMapping(qb.tables)
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

func createJoinMapping(tables map[entities.EntityType]string) map[entities.EntityType]map[entities.EntityType]string {
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

func toAs(field string) string {
	return strings.Replace(field, ".", "_", -1)
}

func removeSchema(table string) string {
	i := strings.Index(table, ".")
	if i == -1 {
		return table
	}

	return table[i+1:]
}

func (qb *queryBuilder) getSelect(et entities.Entity, qo *odata.QueryOptions, addAs bool, selectString string) string {
	var properties []string
	if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
		properties = et.GetPropertyNames()
	} else {
		for _, p := range qo.QuerySelect.Params {
			for _, pn := range et.GetPropertyNames() {
				if strings.ToLower(p) == strings.ToLower(pn) {
					properties = append(properties, pn)
				}
			}
		}
	}

	for _, p := range properties {
		skip := false
		for _, e := range entities.EntityTypeList { //getPropertyNames can contains entities these are not needed in the select string
			if p == e.ToString() {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		toAdd := ""
		if len(selectString) > 0 {
			toAdd += ", "
		}

		field := selectMapping[et.GetEntityType()][p]
		if addAs {
			selectString += fmt.Sprintf("%s%s as %s", toAdd, field, toAs(field))
		} else {
			selectString += fmt.Sprintf("%s%s", toAdd, field)
		}

	}

	if qo != nil && !qo.QueryExpand.IsNil() {
		for _, o := range qo.QueryExpand.Operations {
			selectString = qb.getSelect(o.Entity, o.QueryOptions, addAs, selectString)
		}
	}

	return selectString
}

// CreateQuery creates a new query based on given input
//   e1 = entity to get
//   e2 = from entity
//   id = where e2
// example: Datastreams(1)/Thing = CreateQuery(&entities.Thing, &entities.Datastream, 1, nil)
func (qb *queryBuilder) CreateQuery(e1 entities.Entity, e2 entities.Entity, id interface{}, qo *odata.QueryOptions) (string, error) {
	queryString := ""
	et1 := e1.GetEntityType()
	if e2 == nil { // no second entity given
		queryString = fmt.Sprintf("select %s from %s%s", qb.getSelect(e1, qo, true, ""), qb.tables[et1], qb.createLateralJoin(e1, nil, qo, ""))
		if id != nil {
			queryString = fmt.Sprintf("%s where %s = %v", queryString, selectMapping[et1]["id"], id)
		}
	} else {
		et2 := e2.GetEntityType()
		//queryString = fmt.Sprintf("select %s from %s %s", qb.getSelect(e1, qo, ""), qb.tables[et1], qb.joins[et2][et1])
		queryString = fmt.Sprintf("select %s from %s %s", qb.getSelect(e1, qo, true, ""), qb.tables[et1], qb.createLateralJoin(e1, e2, qo, ""))
		if id != nil {
			queryString = fmt.Sprintf("%s where %s = %v", queryString, selectMapping[et2]["id"], id)
		}
	}

	if qo != nil && !qo.QueryExpand.IsNil() {

	}

	//orderby

	//append limit/offset
	queryString = fmt.Sprintf("%s%s", queryString, CreateTopSkipQueryString(qo))

	return queryString, nil
}

func (qb *queryBuilder) createLateralJoin(e1 entities.Entity, e2 entities.Entity, qo *odata.QueryOptions, joinString string) string {
	if e2 != nil {
		joinString = fmt.Sprintf("%s "+
			"inner join lateral ( "+
			"select %s "+
			"from %s "+
			"where %s"+
			") as %s on true", joinString, qb.getSelect(e2, qo, false, ""), qb.tables[e2.GetEntityType()], qb.joins[e1.GetEntityType()][e2.GetEntityType()], removeSchema(qb.tables[e2.GetEntityType()]))
	} else {
		if qo != nil && !qo.QueryExpand.IsNil() {
			for _, qe := range qo.QueryExpand.Operations {
				joinString = qb.createLateralJoin(e1, qe.Entity, qe.QueryOptions, joinString)
			}
		}
	}

	return joinString
}

/*
SELECT datastream.name AS datastream_name, observation.data AS observation_data, featureofinterest.feature AS FOI_feature
FROM v1.datastream
INNER JOIN LATERAL (
    SELECT observation.data, observation.featureofinterest_id
    FROM v1.observation
    WHERE observation.stream_id = datastream.id
    ORDER BY observation.id DESC
    LIMIT 4
) AS observation
ON TRUE
INNER JOIN LATERAL (
    SELECT observation.data, featureofinterest.feature
    FROM v1.featureofinterest
    WHERE observation.featureofinterest_id = featureofinterest.id
    ORDER BY featureofinterest.id DESC
) AS featureofinterest
ON TRUE
WHERE datastream.id = 2;
*/

func (qb *queryBuilder) ExecuteQuery(query string) {

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
