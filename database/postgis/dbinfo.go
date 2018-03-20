package postgis

import (
	"fmt"

	entities "github.com/gost/core"
	"github.com/gost/godata"
)

// tables as defined in postgis
var (
	thingTable                        = "thing"
	locationTable                     = "location"
	historicalLocationTable           = "historicallocation"
	sensorTable                       = "sensor"
	observedPropertyTable             = "observedproperty"
	datastreamTable                   = "datastream"
	observationTable                  = "observation"
	featureOfInterestTable            = "featureofinterest"
	thingToLocationTable              = "thing_to_location"
	locationToHistoricalLocationTable = "location_to_historicallocation"
)

// thing fields
var (
	thingID          = "id"
	thingName        = "name"
	thingDescription = "description"
	thingProperties  = "properties"
)

// location fields
var (
	locationID           = "id"
	locationName         = "name"
	locationDescription  = "description"
	locationEncodingType = "encodingtype"
	locationLocation     = "location"
	locationGeoJSON      = "geojson"
)

// thingToLocationTable fields
var (
	thingToLocationThingID    = "thing_id"
	thingToLocationLocationID = "location_id"
)

// historical location fields
var (
	historicalLocationID         = "id"
	historicalLocationTime       = "time"
	historicalLocationThingID    = "thing_id"
	historicalLocationLocationID = "location_id"
)

// locationToHistoricalLocation fields
var (
	locationToHistoricalLocationLocationID           = "location_id "
	locationToHistoricalLocationHistoricalLocationID = "historicallocation_id "
)

var (
	asSeparator = "_"
	idField     = "id"
)

// sensor fields
var (
	sensorID           = idField
	sensorName         = "name"
	sensorDescription  = "description"
	sensorEncodingType = "encodingtype"
	sensorMetadata     = "metadata"
)

// observed property fields
var (
	observedPropertyID          = idField
	observedPropertyName        = "name"
	observedPropertyDescription = "description"
	observedPropertyDefinition  = "definition"
)

// datastream fields
var (
	datastreamID                 = idField
	datastreamName               = "name"
	datastreamDescription        = "description"
	datastreamUnitOfMeasurement  = "unitofmeasurement"
	datastreamObservationType    = "observationtype"
	datastreamObservedArea       = "observedarea"
	datastreamPhenomenonTime     = "phenomenontime"
	datastreamResultTime         = "resulttime"
	datastreamThingID            = "thing_id"
	datastreamSensorID           = "sensor_id"
	datastreamObservedPropertyID = "observedproperty_id"
)

// observation fields
var (
	observationID                  = idField
	observationData                = "data"
	observationPhenomenonTime      = "phenomenontime"
	observationResultTime          = "resulttime"
	observationResult              = "result"
	observationValidTime           = "validtime"
	observationResultQuality       = "resultquality"
	observationParameters          = "parameters"
	observationStreamID            = "stream_id"
	observationFeatureOfInterestID = "featureofinterest_id"
)

// feature of interest fields
var (
	foiID                 = idField
	foiName               = "name"
	foiDescription        = "description"
	foiEncodingType       = "encodingtype"
	foiFeature            = "feature"
	foiGeoJSON            = "geojson"
	foiOriginalLocationID = "original_location_id"
)

// ParamFactory receives a map of columns (with select as names) with values an implementation should parse it to the correct entity
type ParamFactory func(values map[string]interface{}) (entities.Entity, error)

var asPrefixArr = [...]string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ",
	"BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ",
}

// QueryParseInfo is constructed based on the input send to the QueryBuilder, with the help of QueryParseInfo
// the response rows from the database can be parsed into the correct entities with their relations and sub entities
type QueryParseInfo struct {
	QueryIndex   int    // Order of quest
	AsPrefix     string // Extra AS that gets added to the join string
	Entity       entities.Entity
	ExpandItem   *godata.ExpandItem
	ParamFactory ParamFactory
	Parent       *QueryParseInfo
	SubEntities  []*QueryParseInfo
}

// Init initialises a QueryParseInfo object by setting al the needed info
func (q *QueryParseInfo) Init(entityType entities.EntityType, queryIndex int, parent *QueryParseInfo, expandItem *godata.ExpandItem) {
	q.Parent = parent
	q.AsPrefix = asPrefixArr[queryIndex]
	q.QueryIndex = queryIndex
	q.ExpandItem = expandItem
	switch e := entityType; e {
	case entities.EntityTypeThing:
		q.Entity = &entities.Thing{}
		q.ParamFactory = thingParamFactory
	case entities.EntityTypeFeatureOfInterest:
		q.Entity = &entities.FeatureOfInterest{}
		q.ParamFactory = featureOfInterestParamFactory
	case entities.EntityTypeLocation:
		q.Entity = &entities.Location{}
		q.ParamFactory = locationParamFactory
	case entities.EntityTypeObservation:
		q.Entity = &entities.Observation{}
		q.ParamFactory = observationParamFactory
	case entities.EntityTypeObservedProperty:
		q.Entity = &entities.ObservedProperty{}
		q.ParamFactory = observedPropertyParamFactory
	case entities.EntityTypeDatastream:
		q.Entity = &entities.Datastream{}
		q.ParamFactory = datastreamParamFactory
	case entities.EntityTypeHistoricalLocation:
		q.Entity = &entities.HistoricalLocation{}
		q.ParamFactory = historicalLocationParamFactory
	case entities.EntityTypeSensor:
		q.Entity = &entities.Sensor{}
		q.ParamFactory = sensorParamFactory
	}
}

// GetParent returns the parent entity for a given list of entity types
func (q *QueryParseInfo) GetParent(etl []entities.EntityType) *QueryParseInfo {
	if len(etl) == 0 {
		return q
	}

	cq := q
	path := ""
	if len(etl) > 1 {
		for i, e := range etl {
			if i+1 == len(etl) {
				continue
			}

			if i == 0 {
				path = fmt.Sprintf("%v", e.ToString())
			} else {
				path = fmt.Sprintf("%v/%v", path, e.ToString())
			}
		}
	}

	if path != "" {
		p := fmt.Sprintf("%s/%s", q.Entity.GetEntityType().ToString(), path)
		for _, se := range q.SubEntities {
			if se.getPath("") == p {
				cq = se
				break
			}
		}
	}

	return cq
}

// QueryInfoExists checks if there is already a QueryParseInfo added based on the entity list
func (q *QueryParseInfo) QueryInfoExists(etl []entities.EntityType) (bool, *QueryParseInfo) {
	path := ""
	for i, e := range etl {
		if i == 0 {
			path = fmt.Sprintf("%v", e.ToString())
		} else {
			path = fmt.Sprintf("%v/%v", path, e.ToString())
		}
	}

	if path != "" {
		for i, se := range q.SubEntities {
			p := fmt.Sprintf("%s/%s", q.Entity.GetEntityType().ToString(), path)
			if se.getPath("") == p {
				return true, q.SubEntities[i]
			}
		}
	}

	return false, nil
}

func (q *QueryParseInfo) getPath(path string) string {
	if len(path) == 0 {
		path = q.Entity.GetEntityType().ToString()
	} else {
		path = fmt.Sprintf("%s/%s", q.Entity.GetEntityType().ToString(), path)
	}

	if q.Parent != nil {
		path = q.Parent.getPath(path)
	}

	return path
}

func (q *QueryParseInfo) getSubEntity(et entities.EntityType) *QueryParseInfo {
	for i, se := range q.SubEntities {
		if et == se.Entity.GetEntityType() {
			return q.SubEntities[i]
		}
	}

	return q
}

// GetQueryParseInfoByQueryIndex returns the QueryParseInfo by a given QueryID, this func should be called from the main
// QueryParseInfo object
func (q *QueryParseInfo) GetQueryParseInfoByQueryIndex(id int) *QueryParseInfo {
	if q.QueryIndex == id {
		return q
	}

	for _, qpi := range q.SubEntities {
		t := qpi.GetQueryParseInfoByQueryIndex(id)
		if t != nil {
			return t
		}
	}

	return nil
}

// GetNextQueryIndex returns the next query index number based on the added entities/sub entities
func (q *QueryParseInfo) GetNextQueryIndex() int {
	qpi := q.GetMainQueryParseInfo()

	qi := qpi.QueryIndex
	if len(qpi.SubEntities) > 0 {
		lastSub := qpi.SubEntities[len(qpi.SubEntities)-1]
		qi = lastSub.getNextQueryIndex() - 1
	}

	return qi + 1
}

// GetNextQueryIndex returns the next query index number based on the added entities/sub entities
func (q *QueryParseInfo) getNextQueryIndex() int {
	qi := q.QueryIndex
	if len(q.SubEntities) > 0 {
		lastSub := q.SubEntities[len(q.SubEntities)-1]
		qi = lastSub.GetNextQueryIndex() - 1
	}

	return qi + 1
}

// GetMainQueryParseInfo returns the first QueryParseInfo from the tree
func (q *QueryParseInfo) GetMainQueryParseInfo() *QueryParseInfo {
	qpi := q
	if qpi.Parent != nil {
		qpi = qpi.Parent.GetMainQueryParseInfo()
	}

	return qpi
}

// GetQueryIDRelationMap returns the query index relations, ie QueryParseInfo with sub entity datastream thing qid = 0, datastream qid = 1
// example: returns [1]0 - datastream (1) relates to thing (0)
// returns nil if no relation exists
func (q *QueryParseInfo) GetQueryIDRelationMap(relationMap map[int]int) map[int]int {
	if relationMap == nil {
		relationMap = map[int]int{}
	}

	if len(q.SubEntities) == 0 {
		return relationMap
	}

	for _, qpi := range q.SubEntities {
		relationMap[qpi.QueryIndex] = qpi.Parent.QueryIndex
		relationMap = qpi.GetQueryIDRelationMap(relationMap)
	}

	return relationMap
}

// Parse receives a map containing row names with their values and executes the set ParamFactory to
// parse the database response into an entity
func (q *QueryParseInfo) Parse(values map[string]interface{}) (entities.Entity, error) {
	return q.ParamFactory(values)
}

var asMappings = map[entities.EntityType]map[string]string{
	entities.EntityTypeThing: {
		thingID:          constructAs(thingTable, thingID),
		thingName:        constructAs(thingTable, thingName),
		thingDescription: constructAs(thingTable, thingDescription),
		thingProperties:  constructAs(thingTable, thingProperties),
	},
	entities.EntityTypeLocation: {
		locationID:           constructAs(locationTable, locationID),
		locationName:         constructAs(locationTable, locationName),
		locationDescription:  constructAs(locationTable, locationDescription),
		locationEncodingType: constructAs(locationTable, locationEncodingType),
		locationLocation:     constructAs(locationTable, locationLocation),
		locationGeoJSON:      constructAs(locationTable, locationGeoJSON),
	},
	entities.EntityTypeThingToLocation: {
		thingToLocationThingID:    constructAs(thingToLocationTable, thingToLocationThingID),
		thingToLocationLocationID: constructAs(thingToLocationTable, thingToLocationLocationID),
	},
	entities.EntityTypeLocationToHistoricalLocation: {
		locationToHistoricalLocationLocationID:           constructAs(locationToHistoricalLocationTable, locationToHistoricalLocationLocationID),
		locationToHistoricalLocationHistoricalLocationID: constructAs(locationToHistoricalLocationTable, locationToHistoricalLocationHistoricalLocationID),
	},
	entities.EntityTypeHistoricalLocation: {
		historicalLocationID:         constructAs(historicalLocationTable, historicalLocationID),
		historicalLocationTime:       constructAs(historicalLocationTable, historicalLocationTime),
		historicalLocationThingID:    constructAs(historicalLocationTable, historicalLocationThingID),
		historicalLocationLocationID: constructAs(historicalLocationTable, historicalLocationLocationID),
	},
	entities.EntityTypeSensor: {
		sensorID:           constructAs(sensorTable, sensorID),
		sensorName:         constructAs(sensorTable, sensorName),
		sensorDescription:  constructAs(sensorTable, sensorDescription),
		sensorEncodingType: constructAs(sensorTable, sensorEncodingType),
		sensorMetadata:     constructAs(sensorTable, sensorMetadata),
	},
	entities.EntityTypeObservedProperty: {
		observedPropertyID:          constructAs(observedPropertyTable, observedPropertyID),
		observedPropertyName:        constructAs(observedPropertyTable, observedPropertyName),
		observedPropertyDescription: constructAs(observedPropertyTable, observedPropertyDescription),
		observedPropertyDefinition:  constructAs(observedPropertyTable, observedPropertyDefinition),
	},
	entities.EntityTypeObservation: {
		observationID:                  constructAs(observationTable, observationID),
		observationData:                constructAs(observationTable, observationData),
		observationPhenomenonTime:      constructAs(observationTable, observationPhenomenonTime),
		observationResultTime:          constructAs(observationTable, observationResultTime),
		observationResult:              constructAs(observationTable, observationResult),
		observationValidTime:           constructAs(observationTable, observationValidTime),
		observationResultQuality:       constructAs(observationTable, observationResultQuality),
		observationParameters:          constructAs(observationTable, observationParameters),
		observationStreamID:            constructAs(observationTable, observationStreamID),
		observationFeatureOfInterestID: constructAs(observationTable, observationFeatureOfInterestID),
	},
	entities.EntityTypeFeatureOfInterest: {
		foiID:                 constructAs(featureOfInterestTable, foiID),
		foiName:               constructAs(featureOfInterestTable, foiName),
		foiDescription:        constructAs(featureOfInterestTable, foiDescription),
		foiEncodingType:       constructAs(featureOfInterestTable, foiEncodingType),
		foiFeature:            constructAs(featureOfInterestTable, foiFeature),
		foiGeoJSON:            constructAs(featureOfInterestTable, foiGeoJSON),
		foiOriginalLocationID: constructAs(featureOfInterestTable, foiOriginalLocationID),
	},
	entities.EntityTypeDatastream: {
		datastreamID:                 constructAs(datastreamTable, datastreamID),
		datastreamName:               constructAs(datastreamTable, datastreamName),
		datastreamDescription:        constructAs(datastreamTable, datastreamDescription),
		datastreamUnitOfMeasurement:  constructAs(datastreamTable, datastreamUnitOfMeasurement),
		datastreamObservationType:    constructAs(datastreamTable, datastreamObservationType),
		datastreamObservedArea:       constructAs(datastreamTable, datastreamObservedArea),
		datastreamPhenomenonTime:     constructAs(datastreamTable, datastreamPhenomenonTime),
		datastreamResultTime:         constructAs(datastreamTable, datastreamResultTime),
		datastreamThingID:            constructAs(datastreamTable, datastreamThingID),
		datastreamSensorID:           constructAs(datastreamTable, datastreamSensorID),
		datastreamObservedPropertyID: constructAs(datastreamTable, datastreamObservedPropertyID),
	},
}

func constructAs(table, field string) string {
	return fmt.Sprintf("%s%s%s", table, asSeparator, field)
}

var tableMappings = map[entities.EntityType]string{
	entities.EntityTypeThing:              thingTable,
	entities.EntityTypeLocation:           locationTable,
	entities.EntityTypeThingToLocation:    thingToLocationTable,
	entities.EntityTypeHistoricalLocation: historicalLocationTable,
	entities.EntityTypeSensor:             sensorTable,
	entities.EntityTypeObservedProperty:   observedPropertyTable,
	entities.EntityTypeObservation:        observationTable,
	entities.EntityTypeFeatureOfInterest:  featureOfInterestTable,
	entities.EntityTypeDatastream:         datastreamTable,
}

var selectAsMappings = map[entities.EntityType]map[string]string{
	entities.EntityTypeThing: {
		thingID: fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeThing], asMappings[entities.EntityTypeThing][thingID]),
	},
	entities.EntityTypeLocation: {
		locationID: fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeLocation], asMappings[entities.EntityTypeLocation][locationID]),
	},
	entities.EntityTypeThingToLocation: {
		thingToLocationThingID:    fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeThingToLocation], asMappings[entities.EntityTypeThingToLocation][thingToLocationThingID]),
		thingToLocationLocationID: fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeThingToLocation], asMappings[entities.EntityTypeThingToLocation][thingToLocationLocationID]),
	},
	entities.EntityTypeHistoricalLocation: {
		historicalLocationID:      fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeHistoricalLocation], asMappings[entities.EntityTypeHistoricalLocation][historicalLocationID]),
		historicalLocationThingID: fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeHistoricalLocation], asMappings[entities.EntityTypeHistoricalLocation][historicalLocationThingID]),
	},
	entities.EntityTypeLocationToHistoricalLocation: {
		locationToHistoricalLocationHistoricalLocationID: fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeLocationToHistoricalLocation], asMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationHistoricalLocationID]),
		locationToHistoricalLocationLocationID:           fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeLocationToHistoricalLocation], asMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationLocationID]),
	},
	entities.EntityTypeSensor: {
		sensorID: fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeSensor], asMappings[entities.EntityTypeSensor][sensorID]),
	},
	entities.EntityTypeObservedProperty: {
		observedPropertyID: fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeObservedProperty], asMappings[entities.EntityTypeObservedProperty][observedPropertyID]),
	},
	entities.EntityTypeObservation: {
		observationID:                  fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeObservation], asMappings[entities.EntityTypeObservation][observationID]),
		observationStreamID:            fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeObservation], asMappings[entities.EntityTypeObservation][observationStreamID]),
		observationFeatureOfInterestID: fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeObservation], asMappings[entities.EntityTypeObservation][observationFeatureOfInterestID]),
	},
	entities.EntityTypeFeatureOfInterest: {
		foiID: fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeFeatureOfInterest], asMappings[entities.EntityTypeFeatureOfInterest][foiID]),
	},
	entities.EntityTypeDatastream: {
		datastreamID:                 fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeDatastream], asMappings[entities.EntityTypeDatastream][datastreamID]),
		datastreamThingID:            fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeDatastream], asMappings[entities.EntityTypeDatastream][datastreamThingID]),
		datastreamObservedPropertyID: fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeDatastream], asMappings[entities.EntityTypeDatastream][datastreamObservedPropertyID]),
		datastreamSensorID:           fmt.Sprintf("%s.%s", tableMappings[entities.EntityTypeDatastream], asMappings[entities.EntityTypeDatastream][datastreamSensorID]),
	},
}

// maps an entity property name to the right field
var selectMappings = map[entities.EntityType]map[string]string{
	entities.EntityTypeThing: {
		thingID:          fmt.Sprintf("%s.%s", thingTable, thingID),
		thingName:        fmt.Sprintf("%s.%s", thingTable, thingName),
		thingDescription: fmt.Sprintf("%s.%s", thingTable, thingDescription),
		thingProperties:  fmt.Sprintf("%s.%s", thingTable, thingProperties),
	},
	entities.EntityTypeLocation: {
		locationID:           fmt.Sprintf("%s.%s", locationTable, locationID),
		locationName:         fmt.Sprintf("%s.%s", locationTable, locationName),
		locationDescription:  fmt.Sprintf("%s.%s", locationTable, locationDescription),
		locationEncodingType: fmt.Sprintf("%s.%s", locationTable, locationEncodingType),
		locationLocation:     fmt.Sprintf("public.ST_AsGeoJSON(%s.%s)", locationTable, locationLocation),
		locationGeoJSON:      fmt.Sprintf("%s.%s::text", locationTable, locationGeoJSON),
	},
	entities.EntityTypeThingToLocation: {
		thingToLocationThingID:    fmt.Sprintf("%s.%s", thingToLocationTable, thingToLocationThingID),
		thingToLocationLocationID: fmt.Sprintf("%s.%s", thingToLocationTable, thingToLocationLocationID),
	},
	entities.EntityTypeHistoricalLocation: {
		historicalLocationID:         fmt.Sprintf("%s.%s", historicalLocationTable, historicalLocationID),
		historicalLocationTime:       fmt.Sprintf("to_char(%s.%s at time zone 'UTC', '%s')", historicalLocationTable, historicalLocationTime, TimeFormat),
		historicalLocationThingID:    fmt.Sprintf("%s.%s", historicalLocationTable, historicalLocationThingID),
		historicalLocationLocationID: fmt.Sprintf("%s.%s", historicalLocationTable, historicalLocationLocationID),
	},
	entities.EntityTypeLocationToHistoricalLocation: {
		locationToHistoricalLocationHistoricalLocationID: fmt.Sprintf("%s.%s", locationToHistoricalLocationTable, locationToHistoricalLocationHistoricalLocationID),
		locationToHistoricalLocationLocationID:           fmt.Sprintf("%s.%s", locationToHistoricalLocationTable, locationToHistoricalLocationLocationID),
	},
	entities.EntityTypeSensor: {
		sensorID:           fmt.Sprintf("%s.%s", sensorTable, sensorID),
		sensorName:         fmt.Sprintf("%s.%s", sensorTable, sensorName),
		sensorDescription:  fmt.Sprintf("%s.%s", sensorTable, sensorDescription),
		sensorEncodingType: fmt.Sprintf("%s.%s", sensorTable, sensorEncodingType),
		sensorMetadata:     fmt.Sprintf("%s.%s", sensorTable, sensorMetadata),
	},
	entities.EntityTypeObservedProperty: {
		observedPropertyID:          fmt.Sprintf("%s.%s", observedPropertyTable, observedPropertyID),
		observedPropertyName:        fmt.Sprintf("%s.%s", observedPropertyTable, observedPropertyName),
		observedPropertyDescription: fmt.Sprintf("%s.%s", observedPropertyTable, observedPropertyDescription),
		observedPropertyDefinition:  fmt.Sprintf("%s.%s", observedPropertyTable, observedPropertyDefinition),
	},
	entities.EntityTypeObservation: {
		observationID:                  fmt.Sprintf("%s.%s", observationTable, observationID),
		observationData:                fmt.Sprintf("%s.%s", observationTable, observationData),
		observationPhenomenonTime:      fmt.Sprintf("%s.%s ->> '%s'", observationTable, observationData, "phenomenonTime"),
		observationResultTime:          fmt.Sprintf("%s.%s ->> '%s'", observationTable, observationData, "resultTime"),
		observationResult:              fmt.Sprintf("%s.%s -> '%s'", observationTable, observationData, observationResult),
		observationValidTime:           fmt.Sprintf("%s.%s ->> '%s'", observationTable, observationData, "validTime"),
		observationResultQuality:       fmt.Sprintf("%s.%s ->> '%s'", observationTable, observationData, "resultQuality"),
		observationParameters:          fmt.Sprintf("%s.%s ->> '%s'", observationTable, observationData, observationParameters),
		observationStreamID:            fmt.Sprintf("%s.%s", observationTable, observationStreamID),
		observationFeatureOfInterestID: fmt.Sprintf("%s.%s", observationTable, observationFeatureOfInterestID),
	},
	entities.EntityTypeFeatureOfInterest: {
		foiID:                 fmt.Sprintf("%s.%s", featureOfInterestTable, foiID),
		foiName:               fmt.Sprintf("%s.%s", featureOfInterestTable, foiName),
		foiDescription:        fmt.Sprintf("%s.%s", featureOfInterestTable, foiDescription),
		foiEncodingType:       fmt.Sprintf("%s.%s", featureOfInterestTable, foiEncodingType),
		foiFeature:            fmt.Sprintf("public.ST_AsGeoJSON(%s.%s)", featureOfInterestTable, foiFeature),
		foiGeoJSON:            fmt.Sprintf("%s.%s::text", featureOfInterestTable, foiGeoJSON),
		foiOriginalLocationID: fmt.Sprintf("%s.%s", featureOfInterestTable, foiOriginalLocationID),
	},
	entities.EntityTypeDatastream: {
		datastreamID:                 fmt.Sprintf("%s.%s", datastreamTable, datastreamID),
		datastreamName:               fmt.Sprintf("%s.%s", datastreamTable, datastreamName),
		datastreamDescription:        fmt.Sprintf("%s.%s", datastreamTable, datastreamDescription),
		datastreamUnitOfMeasurement:  fmt.Sprintf("%s.%s", datastreamTable, datastreamUnitOfMeasurement),
		datastreamObservationType:    fmt.Sprintf("%s.%s", datastreamTable, datastreamObservationType),
		datastreamObservedArea:       fmt.Sprintf("public.ST_AsGeoJSON(%s.%s)", datastreamTable, datastreamObservedArea),
		datastreamPhenomenonTime:     fmt.Sprintf("%s.%s", datastreamTable, datastreamPhenomenonTime),
		datastreamResultTime:         fmt.Sprintf("%s.%s", datastreamTable, datastreamResultTime),
		datastreamThingID:            fmt.Sprintf("%s.%s", datastreamTable, datastreamThingID),
		datastreamSensorID:           fmt.Sprintf("%s.%s", datastreamTable, datastreamSensorID),
		datastreamObservedPropertyID: fmt.Sprintf("%s.%s", datastreamTable, datastreamObservedPropertyID),
	},
}

var selectMappingsIgnore = map[entities.EntityType]map[string]bool{
	entities.EntityTypeLocation: {
		locationLocation: true,
	},
}

func getJoin(tableMap map[entities.EntityType]string, get entities.EntityType, by entities.EntityType, asPrefix string) string {
	switch get {
	case entities.EntityTypeThing:
		{
			return getJoinThing(tableMap, by, asPrefix)
		}
	case entities.EntityTypeLocation: // get Location by ...
		{
			return getJoinLocation(tableMap, by, asPrefix)

		}
	case entities.EntityTypeHistoricalLocation: // get HistoricalLocation by ... //fmt.Sprintf("%s WHERE %s.%s = %v", queryString, tableMappings[et2], asMappings[et2][idField], id)
		{
			return getJoinHistoricalLocation(tableMap, by, asPrefix)

		}
	case entities.EntityTypeSensor: // get sensor by ...
		{
			return getJoinSensor(tableMap, by, asPrefix)

		}
	case entities.EntityTypeObservedProperty: // get observed property by ...
		{
			return getJoinObservedProperty(tableMap, by, asPrefix)

		}
	case entities.EntityTypeObservation: // get observation by ...
		{
			return getJoinObservations(tableMap, by, asPrefix)

		}
	case entities.EntityTypeFeatureOfInterest: // get feature of interest by ...
		{
			return getJoinFeatureOfInterest(tableMap, by, asPrefix)

		}
	case entities.EntityTypeDatastream: // get Datastream by ...
		{
			return getJoinDatastream(tableMap, by, asPrefix)
		}
	}

	return ""
}

func getJoinThing(tableMap map[entities.EntityType]string, by entities.EntityType, asPrefix string) string {
	switch by {
	case entities.EntityTypeDatastream:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeThing][thingID], createWhereIs(entities.EntityTypeDatastream, datastreamThingID, asPrefix))
	case entities.EntityTypeHistoricalLocation:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeThing][thingID], createWhereIs(entities.EntityTypeHistoricalLocation, historicalLocationThingID, asPrefix))
	case entities.EntityTypeLocation:
		return fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
			tableMap[entities.EntityTypeThingToLocation],
			selectMappings[entities.EntityTypeThing][thingID],
			selectMappings[entities.EntityTypeThingToLocation][thingToLocationThingID],
			createWhereIs(entities.EntityTypeLocation, thingID, asPrefix),
			selectMappings[entities.EntityTypeThingToLocation][thingToLocationLocationID])
	}

	return ""
}

func getJoinLocation(tableMap map[entities.EntityType]string, by entities.EntityType, asPrefix string) string {
	switch by {
	case entities.EntityTypeHistoricalLocation:
		return fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
			tableMap[entities.EntityTypeLocationToHistoricalLocation],
			selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationLocationID],
			selectMappings[entities.EntityTypeLocation][locationID],
			selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationHistoricalLocationID],
			createWhereIs(entities.EntityTypeHistoricalLocation, historicalLocationID, asPrefix))
	case entities.EntityTypeThing:
		return fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
			tableMap[entities.EntityTypeThingToLocation],
			selectMappings[entities.EntityTypeThingToLocation][thingToLocationLocationID],
			selectMappings[entities.EntityTypeLocation][locationID],
			selectMappings[entities.EntityTypeThingToLocation][thingToLocationThingID],
			createWhereIs(entities.EntityTypeThing, thingID, asPrefix))
	}

	return ""
}

func getJoinHistoricalLocation(tableMap map[entities.EntityType]string, by entities.EntityType, asPrefix string) string {
	switch by {
	case entities.EntityTypeLocation:
		return fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
			tableMap[entities.EntityTypeLocationToHistoricalLocation],
			selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationHistoricalLocationID],
			selectMappings[entities.EntityTypeHistoricalLocation][historicalLocationID],
			selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationLocationID],
			createWhereIs(entities.EntityTypeLocation, locationID, asPrefix))
	case entities.EntityTypeThing:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeHistoricalLocation][historicalLocationThingID], createWhereIs(entities.EntityTypeThing, thingID, asPrefix))
	}

	return ""
}

func getJoinDatastream(tableMap map[entities.EntityType]string, by entities.EntityType, asPrefix string) string {
	switch by {
	case entities.EntityTypeThing:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeDatastream][datastreamThingID], createWhereIs(entities.EntityTypeThing, thingID, asPrefix))
	case entities.EntityTypeSensor:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeDatastream][datastreamSensorID], createWhereIs(entities.EntityTypeSensor, sensorID, asPrefix))
	case entities.EntityTypeObservedProperty:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeDatastream][datastreamObservedPropertyID], createWhereIs(entities.EntityTypeObservedProperty, observedPropertyID, asPrefix))
	case entities.EntityTypeObservation:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeDatastream][datastreamID], createWhereIs(entities.EntityTypeObservation, observationStreamID, asPrefix))
	case entities.EntityTypeLocation:
		return fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
			tableMap[entities.EntityTypeThingToLocation],
			createWhereIs(entities.EntityTypeLocation, locationID, asPrefix),
			selectMappings[entities.EntityTypeThingToLocation][thingToLocationLocationID],
			selectMappings[entities.EntityTypeThingToLocation][thingToLocationThingID],
			selectMappings[entities.EntityTypeDatastream][datastreamThingID],
		)
	}

	return ""
}

func getJoinSensor(tableMap map[entities.EntityType]string, by entities.EntityType, asPrefix string) string {
	switch by {
	case entities.EntityTypeDatastream:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeSensor][sensorID], createWhereIs(entities.EntityTypeDatastream, datastreamSensorID, asPrefix))
	}

	return ""
}

func getJoinObservedProperty(tableMap map[entities.EntityType]string, by entities.EntityType, asPrefix string) string {
	switch by {
	case entities.EntityTypeDatastream:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeObservedProperty][observedPropertyID], createWhereIs(entities.EntityTypeDatastream, datastreamObservedPropertyID, asPrefix))
	}

	return ""
}

func getJoinObservations(tableMap map[entities.EntityType]string, by entities.EntityType, asPrefix string) string {
	switch by {
	case entities.EntityTypeDatastream:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeObservation][observationStreamID], createWhereIs(entities.EntityTypeDatastream, datastreamID, asPrefix))
	case entities.EntityTypeFeatureOfInterest:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeObservation][observationFeatureOfInterestID], createWhereIs(entities.EntityTypeFeatureOfInterest, foiID, asPrefix))
	}

	return ""
}

func getJoinFeatureOfInterest(tableMap map[entities.EntityType]string, by entities.EntityType, asPrefix string) string {
	switch by {
	case entities.EntityTypeObservation:
		return fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeFeatureOfInterest][foiID], createWhereIs(entities.EntityTypeObservation, observationFeatureOfInterestID, asPrefix))
	}

	return ""
}

func getJoinByID(tableMap map[entities.EntityType]string, get entities.EntityType, by entities.EntityType, id interface{}) string {
	switch get {

	case entities.EntityTypeHistoricalLocation: // get HistoricalLocation by ... //fmt.Sprintf("%s WHERE %s.%s = %v", queryString, tableMappings[et2], asMappings[et2][idField], id)
		{
			return getJoinHistoricalLocationByID(tableMap, by, id)

		}
	case entities.EntityTypeObservation: // get observation by ...
		{
			return getJoinObservationsByID(tableMap, by, id)

		}
	case entities.EntityTypeDatastream: // get Datastream by ...
		{
			return getJoinDatastreamByID(tableMap, by, id)
		}
	}

	return ""
}

func getJoinHistoricalLocationByID(tableMap map[entities.EntityType]string, by entities.EntityType, id interface{}) string {
	switch by {
	case entities.EntityTypeThing:
		return fmt.Sprintf("%s = %v", selectMappings[entities.EntityTypeHistoricalLocation][historicalLocationThingID], id)
	}

	return ""
}

// Sensor(1)/Datastream
func getJoinDatastreamByID(tableMap map[entities.EntityType]string, by entities.EntityType, id interface{}) string {
	switch by {
	case entities.EntityTypeThing:
		return fmt.Sprintf("%s = %v", selectMappings[entities.EntityTypeDatastream][datastreamThingID], id)
	case entities.EntityTypeSensor:
		return fmt.Sprintf("%s = %v", selectMappings[entities.EntityTypeDatastream][datastreamSensorID], id)
	case entities.EntityTypeObservedProperty:
		return fmt.Sprintf("%s = %v", selectMappings[entities.EntityTypeDatastream][datastreamObservedPropertyID], id)
	}

	return ""
}

// Datastreams(1)/Observations
func getJoinObservationsByID(tableMap map[entities.EntityType]string, by entities.EntityType, id interface{}) string {
	switch by {
	case entities.EntityTypeDatastream:
		return fmt.Sprintf("%s = %v", selectMappings[entities.EntityTypeObservation][observationStreamID], id)
	case entities.EntityTypeFeatureOfInterest:
		return fmt.Sprintf("%s = %v", selectMappings[entities.EntityTypeObservation][observationFeatureOfInterestID], id)
	}

	return ""
}

func createWhereIs(et entities.EntityType, field string, asPrefix string) string {
	if len(asPrefix) == 0 {
		return selectMappings[et][field]
	}

	return fmt.Sprintf("%v_%v", asPrefix, selectAsMappings[et][field])
}

func createTableMappings(schema string) map[entities.EntityType]string {
	if len(schema) > 0 {
		schema = fmt.Sprintf("%s.", schema)
	}

	tables := map[entities.EntityType]string{
		entities.EntityTypeThing:                        fmt.Sprintf("%s%s", schema, thingTable),
		entities.EntityTypeLocation:                     fmt.Sprintf("%s%s", schema, locationTable),
		entities.EntityTypeHistoricalLocation:           fmt.Sprintf("%s%s", schema, historicalLocationTable),
		entities.EntityTypeSensor:                       fmt.Sprintf("%s%s", schema, sensorTable),
		entities.EntityTypeObservedProperty:             fmt.Sprintf("%s%s", schema, observedPropertyTable),
		entities.EntityTypeDatastream:                   fmt.Sprintf("%s%s", schema, datastreamTable),
		entities.EntityTypeObservation:                  fmt.Sprintf("%s%s", schema, observationTable),
		entities.EntityTypeFeatureOfInterest:            fmt.Sprintf("%s%s", schema, featureOfInterestTable),
		entities.EntityTypeThingToLocation:              fmt.Sprintf("%s%s", schema, thingToLocationTable),
		entities.EntityTypeLocationToHistoricalLocation: fmt.Sprintf("%s%s", schema, locationToHistoricalLocationTable),
	}

	return tables
}
