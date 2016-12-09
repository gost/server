package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
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

// sensor fields
var (
	sensorID           = "id"
	sensorName         = "name"
	sensorDescription  = "description"
	sensorEncodingType = "encodingtype"
	sensorMetadata     = "metadata"
)

// observed property fields
var (
	observedPropertyID          = "id"
	observedPropertyName        = "name"
	observedPropertyDescription = "description"
	observedPropertyDefinition  = "definition"
)

// datastream fields
var (
	datastreamID                 = "id"
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
	observationID                  = "id"
	observationData                = "data"
	observationPhenomenonTime      = "phenomenonTime"
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
	foiID                 = "id"
	foiName               = "name"
	foiDescription        = "description"
	foiEncodingType       = "encodingtype"
	foiFeature            = "feature"
	foiOriginalLocationID = "original_location_id"
)

type ParamFactory func(values map[string]interface{}) (entities.Entity, error)

type QueryParseInfo struct {
	QueryIndex   int
	ParamFactory ParamFactory
	EntityType   entities.EntityType
	Entity       entities.Entity
	SubEntities  []QueryParseInfo
}

func (q *QueryParseInfo) Init(entityType entities.EntityType, queryIndex int, paramFactory ParamFactory) {
	q.QueryIndex = queryIndex
	q.ParamFactory = paramFactory
	q.EntityType = entityType
	switch e := entityType; e {
	case entities.EntityTypeThing:
		q.Entity = &entities.Thing{}
		break
	case entities.EntityTypeFeatureOfInterest:
		q.Entity = &entities.FeatureOfInterest{}
		break
	case entities.EntityTypeLocation:
		q.Entity = &entities.Location{}
		break
	case entities.EntityTypeObservation:
		q.Entity = &entities.Observation{}
		break
	case entities.EntityTypeObservedProperty:
		q.Entity = &entities.ObservedProperty{}
		break
	case entities.EntityTypeDatastream:
		q.Entity = &entities.Datastream{}
		break
	case entities.EntityTypeHistoricalLocation:
		q.Entity = &entities.HistoricalLocation{}
		break
	case entities.EntityTypeSensor:
		q.Entity = &entities.Sensor{}
		break
	}
}

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

func (q *QueryParseInfo) Parse(values map[string]interface{}) (entities.Entity, error) {
	return q.ParamFactory(values)
}

var asMappings = map[entities.EntityType]map[string]string{
	entities.EntityTypeThing: {
		thingID:          fmt.Sprintf("%s_%s", thingTable, thingID),
		thingName:        fmt.Sprintf("%s_%s", thingTable, thingName),
		thingDescription: fmt.Sprintf("%s_%s", thingTable, thingDescription),
		thingProperties:  fmt.Sprintf("%s_%s", thingTable, thingProperties),
	},
	entities.EntityTypeLocation: {
		locationID:           fmt.Sprintf("%s_%s", locationTable, locationID),
		locationName:         fmt.Sprintf("%s_%s", locationTable, locationName),
		locationDescription:  fmt.Sprintf("%s_%s", locationTable, locationDescription),
		locationEncodingType: fmt.Sprintf("%s_%s", locationTable, locationEncodingType),
		locationLocation:     fmt.Sprintf("%s_%s", locationTable, locationLocation),
	},
	entities.EntityTypeThingToLocation: {
		thingToLocationThingID:    fmt.Sprintf("%s_%s", thingToLocationTable, thingToLocationThingID),
		thingToLocationLocationID: fmt.Sprintf("%s_%s", thingToLocationTable, thingToLocationLocationID),
	},
	entities.EntityTypeLocationToHistoricalLocation: {
		locationToHistoricalLocationLocationID:           fmt.Sprintf("%s_%s", locationToHistoricalLocationTable, locationToHistoricalLocationLocationID),
		locationToHistoricalLocationHistoricalLocationID: fmt.Sprintf("%s_%s", locationToHistoricalLocationTable, locationToHistoricalLocationHistoricalLocationID),
	},
	entities.EntityTypeHistoricalLocation: {
		historicalLocationID:         fmt.Sprintf("%s_%s", historicalLocationTable, historicalLocationID),
		historicalLocationTime:       fmt.Sprintf("%s_%s", historicalLocationTable, historicalLocationTime),
		historicalLocationThingID:    fmt.Sprintf("%s_%s", historicalLocationTable, historicalLocationThingID),
		historicalLocationLocationID: fmt.Sprintf("%s_%s", historicalLocationTable, historicalLocationLocationID),
	},
	entities.EntityTypeSensor: {
		sensorID:           fmt.Sprintf("%s_%s", sensorTable, sensorID),
		sensorName:         fmt.Sprintf("%s_%s", sensorTable, sensorName),
		sensorDescription:  fmt.Sprintf("%s_%s", sensorTable, sensorDescription),
		sensorEncodingType: fmt.Sprintf("%s_%s", sensorTable, sensorEncodingType),
		sensorMetadata:     fmt.Sprintf("%s_%s", sensorTable, sensorMetadata),
	},
	entities.EntityTypeObservedProperty: {
		observedPropertyID:          fmt.Sprintf("%s_%s", observedPropertyTable, observedPropertyID),
		observedPropertyName:        fmt.Sprintf("%s_%s", observedPropertyTable, observedPropertyName),
		observedPropertyDescription: fmt.Sprintf("%s_%s", observedPropertyTable, observedPropertyDescription),
		observedPropertyDefinition:  fmt.Sprintf("%s_%s", observedPropertyTable, observedPropertyDefinition),
	},
	entities.EntityTypeObservation: {
		observationID:                  fmt.Sprintf("%s_%s", observationTable, observationID),
		observationData:                fmt.Sprintf("%s_%s", observationTable, observationData),
		observationPhenomenonTime:      fmt.Sprintf("%s_%s", observationTable, observationPhenomenonTime),
		observationResultTime:          fmt.Sprintf("%s_%s", observationTable, observationResultTime),
		observationResult:              fmt.Sprintf("%s_%s", observationTable, observationResult),
		observationValidTime:           fmt.Sprintf("%s_%s", observationTable, observationValidTime),
		observationResultQuality:       fmt.Sprintf("%s_%s", observationTable, observationResultQuality),
		observationParameters:          fmt.Sprintf("%s_%s", observationTable, observationParameters),
		observationStreamID:            fmt.Sprintf("%s_%s", observationTable, observationStreamID),
		observationFeatureOfInterestID: fmt.Sprintf("%s_%s", observationTable, observationFeatureOfInterestID),
	},
	entities.EntityTypeFeatureOfInterest: {
		foiID:                 fmt.Sprintf("%s_%s", featureOfInterestTable, foiID),
		foiName:               fmt.Sprintf("%s_%s", featureOfInterestTable, foiName),
		foiDescription:        fmt.Sprintf("%s_%s", featureOfInterestTable, foiDescription),
		foiEncodingType:       fmt.Sprintf("%s_%s", featureOfInterestTable, foiEncodingType),
		foiFeature:            fmt.Sprintf("%s_%s)", featureOfInterestTable, foiFeature),
		foiOriginalLocationID: fmt.Sprintf("%s_%s", featureOfInterestTable, foiOriginalLocationID),
	},
	entities.EntityTypeDatastream: {
		datastreamID:                 fmt.Sprintf("%s_%s", datastreamTable, datastreamID),
		datastreamName:               fmt.Sprintf("%s_%s", datastreamTable, datastreamName),
		datastreamDescription:        fmt.Sprintf("%s_%s", datastreamTable, datastreamDescription),
		datastreamUnitOfMeasurement:  fmt.Sprintf("%s_%s", datastreamTable, datastreamUnitOfMeasurement),
		datastreamObservationType:    fmt.Sprintf("%s_%s", datastreamTable, datastreamObservationType),
		datastreamObservedArea:       fmt.Sprintf("%s_%s", datastreamTable, datastreamObservedArea),
		datastreamPhenomenonTime:     fmt.Sprintf("%s_%s", datastreamTable, datastreamPhenomenonTime),
		datastreamResultTime:         fmt.Sprintf("%s_%s", datastreamTable, datastreamResultTime),
		datastreamThingID:            fmt.Sprintf("%s_%s", datastreamTable, datastreamThingID),
		datastreamSensorID:           fmt.Sprintf("%s_%s", datastreamTable, datastreamSensorID),
		datastreamObservedPropertyID: fmt.Sprintf("%s_%s", datastreamTable, datastreamObservedPropertyID),
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
	},
	entities.EntityTypeThingToLocation: {
		thingToLocationThingID:    fmt.Sprintf("%s.%s", thingToLocationTable, thingToLocationThingID),
		thingToLocationLocationID: fmt.Sprintf("%s.%s", thingToLocationTable, thingToLocationLocationID),
	},
	entities.EntityTypeHistoricalLocation: {
		historicalLocationID:         fmt.Sprintf("%s.%s", historicalLocationTable, historicalLocationID),
		historicalLocationTime:       fmt.Sprintf("%s.%s", historicalLocationTable, historicalLocationTime),
		historicalLocationThingID:    fmt.Sprintf("%s.%s", historicalLocationTable, historicalLocationThingID),
		historicalLocationLocationID: fmt.Sprintf("%s.%s", historicalLocationTable, historicalLocationLocationID),
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
		observationPhenomenonTime:      "data -> 'phenomenonTime'",
		observationResultTime:          "data -> 'resultTime'",
		observationResult:              "data -> 'result'",
		observationValidTime:           "data -> 'validTime'",
		observationResultQuality:       "data -> 'resultQuality'",
		observationParameters:          "data -> 'parameters'",
		observationStreamID:            fmt.Sprintf("%s.%s", observationTable, observationStreamID),
		observationFeatureOfInterestID: fmt.Sprintf("%s.%s", observationTable, observationFeatureOfInterestID),
	},
	entities.EntityTypeFeatureOfInterest: {
		foiID:                 fmt.Sprintf("%s.%s", featureOfInterestTable, foiID),
		foiName:               fmt.Sprintf("%s.%s", featureOfInterestTable, foiName),
		foiDescription:        fmt.Sprintf("%s.%s", featureOfInterestTable, foiDescription),
		foiEncodingType:       fmt.Sprintf("%s.%s", featureOfInterestTable, foiEncodingType),
		foiFeature:            fmt.Sprintf("public.ST_AsGeoJSON(%s.%s)", featureOfInterestTable, foiFeature),
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

func createJoinMappings(tableMappings map[entities.EntityType]string) map[entities.EntityType]map[entities.EntityType]string {
	joinMappings := map[entities.EntityType]map[entities.EntityType]string{
		entities.EntityTypeThing: { // get thing by ...
			entities.EntityTypeDatastream:         fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeThing][thingID], selectMappings[entities.EntityTypeDatastream][datastreamThingID]),
			entities.EntityTypeHistoricalLocation: fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeThing][thingID], selectMappings[entities.EntityTypeHistoricalLocation][historicalLocationThingID]),
			entities.EntityTypeLocation: fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
				tableMappings[entities.EntityTypeThingToLocation],
				selectMappings[entities.EntityTypeThing][thingID],
				selectMappings[entities.EntityTypeThingToLocation][thingToLocationThingID],
				selectMappings[entities.EntityTypeLocation][thingID],
				selectMappings[entities.EntityTypeThingToLocation][thingToLocationLocationID]),
		},
		entities.EntityTypeLocation: { // get Location by ...
			entities.EntityTypeHistoricalLocation: fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
				tableMappings[entities.EntityTypeLocationToHistoricalLocation],
				selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationLocationID],
				selectMappings[entities.EntityTypeLocation][locationID],
				selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationHistoricalLocationID],
				selectMappings[entities.EntityTypeHistoricalLocation][historicalLocationID]),
			entities.EntityTypeThing: fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
				tableMappings[entities.EntityTypeThingToLocation],
				selectMappings[entities.EntityTypeThingToLocation][thingToLocationLocationID],
				selectMappings[entities.EntityTypeLocation][locationID],
				selectMappings[entities.EntityTypeThingToLocation][thingToLocationThingID],
				selectMappings[entities.EntityTypeThing][thingID]),
		},
		entities.EntityTypeHistoricalLocation: { // get HistoricalLocation by ...
			entities.EntityTypeLocation: fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
				tableMappings[entities.EntityTypeLocationToHistoricalLocation],
				selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationHistoricalLocationID],
				selectMappings[entities.EntityTypeHistoricalLocation][historicalLocationID],
				selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationLocationID],
				selectMappings[entities.EntityTypeLocation][locationID]),
			entities.EntityTypeThing: fmt.Sprintf("WHERE %s = %s", historicalLocationThingID, thingID),
		},
		entities.EntityTypeSensor: { // get sensor by ...
			entities.EntityTypeDatastream: fmt.Sprintf("WHERE %s = %s", sensorID, datastreamSensorID),
		},
		entities.EntityTypeObservedProperty: { // get observed property by ...
			entities.EntityTypeDatastream: fmt.Sprintf("WHERE %s = %s", observedPropertyID, datastreamObservedPropertyID),
		},
		entities.EntityTypeObservation: { // get observation by ...
			entities.EntityTypeDatastream:        fmt.Sprintf("WHERE %s = %s", observationStreamID, datastreamID),
			entities.EntityTypeFeatureOfInterest: fmt.Sprintf("WHERE %s = %s", observationFeatureOfInterestID, foiID),
		},
		entities.EntityTypeFeatureOfInterest: { // get feature of interest by ...
			entities.EntityTypeObservation: fmt.Sprintf("WHERE %s = %s", foiID, observationFeatureOfInterestID),
		},
		entities.EntityTypeDatastream: { // get Datastream by ...
			entities.EntityTypeThing:            fmt.Sprintf("WHERE %s = %s", datastreamThingID, thingID),
			entities.EntityTypeSensor:           fmt.Sprintf("WHERE %s = %s", datastreamSensorID, sensorID),
			entities.EntityTypeObservedProperty: fmt.Sprintf("WHERE %s = %s", datastreamObservedPropertyID, observedPropertyID),
			entities.EntityTypeObservation:      fmt.Sprintf("WHERE %s = %s", datastreamID, observationStreamID),
		},
	}

	return joinMappings
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
