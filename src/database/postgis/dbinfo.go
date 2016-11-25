package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
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
	thingToLocationTable    = "thing_to_location"
)

// thing fields
var (
	thingID          = fmt.Sprintf("%s.id", thingTable)
	thingName        = fmt.Sprintf("%s.name", thingTable)
	thingDescription = fmt.Sprintf("%s.description", thingTable)
	thingProperties  = fmt.Sprintf("%s.properties", thingTable)
)

// location fields
var (
	locationID           = fmt.Sprintf("%s.id", locationTable)
	locationName         = fmt.Sprintf("%s.name", locationTable)
	locationDescription  = fmt.Sprintf("%s.description", locationTable)
	locationEncodingType = fmt.Sprintf("%s.encodingtype", locationTable)
	locationLocation     = fmt.Sprintf("%s.location", locationTable)
)

// thingToLocationTable fields
var (
	thingToLocationThingID    = fmt.Sprintf("%s.thing_id", thingToLocationTable)
	thingToLocationLocationID = fmt.Sprintf("%s.location_id", thingToLocationTable)
)

// historical location fields
var (
	historicalLocationID         = fmt.Sprintf("%s.id", historicalLocationTable)
	historicalLocationTime       = fmt.Sprintf("%s.time", historicalLocationTable)
	historicalLocationThingID    = fmt.Sprintf("%s.thing_id", historicalLocationTable)
	historicalLocationLocationID = fmt.Sprintf("%s.location_id", historicalLocationTable)
)

// sensor fields
var (
	sensorID           = fmt.Sprintf("%s.id", sensorTable)
	sensorName         = fmt.Sprintf("%s.name", sensorTable)
	sensorDescription  = fmt.Sprintf("%s.description", sensorTable)
	sensorEncodingType = fmt.Sprintf("%s.encodingtype", sensorTable)
	sensorMetadata     = fmt.Sprintf("%s.metadata", sensorTable)
)

// observed property fields
var (
	observedPropertyID          = fmt.Sprintf("%s.id", observedPropertyTable)
	observedPropertyName        = fmt.Sprintf("%s.name", observedPropertyTable)
	observedPropertyDescription = fmt.Sprintf("%s.description", observedPropertyTable)
	observedPropertyDefinition  = fmt.Sprintf("%s.definition", observedPropertyTable)
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

// observation fields
var (
	observationID                  = fmt.Sprintf("%s.id", observationTable)
	observationData                = fmt.Sprintf("%s.data", observationTable)
	observationStreamID            = fmt.Sprintf("%s.stream_id", observationTable)
	observationFeatureOfInterestID = fmt.Sprintf("%s.featureofinterest_id", observationTable)
)

// feature of interest fields
var (
	foiID                 = fmt.Sprintf("%s.id", featureOfInterestTable)
	foiName               = fmt.Sprintf("%s.name", featureOfInterestTable)
	foiDescription        = fmt.Sprintf("%s.description", featureOfInterestTable)
	foiEncodingType       = fmt.Sprintf("%s.encodingtype", featureOfInterestTable)
	foiFeature            = fmt.Sprintf("%s.feature", featureOfInterestTable)
	foiOriginalLocationID = fmt.Sprintf("%s.original_location_id", featureOfInterestTable)
)

// maps an entity property name to the right field
var selectMappings = map[entities.EntityType]map[string]string{
	entities.EntityTypeThing: {
		"id":          thingID,
		"name":        thingName,
		"description": thingDescription,
		"properties":  thingProperties,
	},
	entities.EntityTypeLocation: {
		"id":           locationID,
		"name":         locationName,
		"description":  locationDescription,
		"encodingtype": locationEncodingType,
		"location":     fmt.Sprintf("public.ST_AsGeoJSON(%s)", locationLocation),
	},
	entities.EntityTypeThingToLocation: {
		"thing_id":    thingToLocationThingID,
		"location_id": thingToLocationLocationID,
	},
	entities.EntityTypeHistoricalLocation: {
		"id":          historicalLocationID,
		"time":        historicalLocationTime,
		"thing_id":    historicalLocationThingID,
		"location_id": historicalLocationLocationID,
	},
	entities.EntityTypeSensor: {
		"id":           sensorID,
		"name":         sensorName,
		"description":  sensorDescription,
		"encodingtype": sensorEncodingType,
		"metadata":     sensorMetadata,
	},
	entities.EntityTypeObservedProperty: {
		"id":          observedPropertyID,
		"name":        observedPropertyName,
		"description": observedPropertyDescription,
		"definition":  observedPropertyDefinition,
	},
	entities.EntityTypeObservation: {
		"id":                   observationID,
		"data":                 observationData,
		"phenomenontime":       "data -> 'phenomenonTime'",
		"resulttime":           "data -> 'resultTime'",
		"result":               "data -> 'result'",
		"validtime":            "data -> 'validTime'",
		"resultquality":        "data -> 'resultQuality'",
		"parameters":           "data -> 'parameters'",
		"stream_id":            observationStreamID,
		"featureofinterest_id": observationFeatureOfInterestID,
	},
	entities.EntityTypeFeatureOfInterest: {
		"id":                   foiID,
		"name":                 foiName,
		"description":          foiDescription,
		"encodingtype":         foiEncodingType,
		"feature":              fmt.Sprintf("public.ST_AsGeoJSON(%s)", foiFeature),
		"original_location_id": foiOriginalLocationID,
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

func createJoinMappings(tableMappings map[entities.EntityType]string) map[entities.EntityType]map[entities.EntityType]string {
	joinMappings := map[entities.EntityType]map[entities.EntityType]string{
		entities.EntityTypeThing: { // get thing by ...
			entities.EntityTypeDatastream:         fmt.Sprintf("WHERE %s = %s", thingID, datastreamThingID),
			entities.EntityTypeHistoricalLocation: fmt.Sprintf("WHERE %s = %s", thingID, historicalLocationThingID),
			entities.EntityTypeLocation: fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
				tableMappings[entities.EntityTypeThingToLocation],
				selectMappings[entities.EntityTypeThing]["id"],
				selectMappings[entities.EntityTypeThingToLocation]["thing_id"],
				selectMappings[entities.EntityTypeLocation]["id"],
				selectMappings[entities.EntityTypeThingToLocation]["location_id"]),
		},
		entities.EntityTypeLocation: { // get Location by ...
			entities.EntityTypeHistoricalLocation: fmt.Sprintf("WHERE %s = %s", locationID, historicalLocationLocationID),
			entities.EntityTypeThing: fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
				tableMappings[entities.EntityTypeThingToLocation],
				selectMappings[entities.EntityTypeThingToLocation]["location_id"],
				selectMappings[entities.EntityTypeLocation]["id"],
				selectMappings[entities.EntityTypeThingToLocation]["thing_id"],
				selectMappings[entities.EntityTypeThing]["id"]),
		},
		entities.EntityTypeHistoricalLocation: { // get Location by ...
			entities.EntityTypeLocation: fmt.Sprintf("WHERE %s = %s", historicalLocationLocationID, locationID),
			entities.EntityTypeThing:    fmt.Sprintf("WHERE %s = %s", historicalLocationThingID, thingID),
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
		entities.EntityTypeThing:              fmt.Sprintf("%s%s", schema, thingTable),
		entities.EntityTypeLocation:           fmt.Sprintf("%s%s", schema, locationTable),
		entities.EntityTypeHistoricalLocation: fmt.Sprintf("%s%s", schema, historicalLocationTable),
		entities.EntityTypeSensor:             fmt.Sprintf("%s%s", schema, sensorTable),
		entities.EntityTypeObservedProperty:   fmt.Sprintf("%s%s", schema, observedPropertyTable),
		entities.EntityTypeDatastream:         fmt.Sprintf("%s%s", schema, datastreamTable),
		entities.EntityTypeObservation:        fmt.Sprintf("%s%s", schema, observationTable),
		entities.EntityTypeFeatureOfInterest:  fmt.Sprintf("%s%s", schema, featureOfInterestTable),
		entities.EntityTypeThingToLocation:    fmt.Sprintf("%s%s", schema, thingToLocationTable),
	}

	return tables
}
