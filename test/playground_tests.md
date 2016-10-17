## Playground tests

Here a collection of sample requests to start working with the GOST API. It contains 6 steps: create a Thing, create a Location, create an ObservedProperty, create a Sensor, create a Datastream and create an Observation.

Instructions:

. replace http://localhost:8080 with the url of your own test server

. when something is created with a HTTP Post request, an @iot.id is returned. Use this @iot.id in subsequent requests (see parameters)

## 1] Create a Thing

A thing is an object in either the physical or virtual world.

Request:

```sh
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "description": "my thermometer",
    "name":"BertThermometer",
    "properties": {
        "organisation": "Geodan",
        "owner": "Bert"
    }
}' "http://localhost:8080/v1.0/Things"
```

Parameters: -

Response:
```json
{
   "@iot.id": "5",
   "@iot.selfLink": "http://localhost:8080/v1.0/Things(5)",
   "description": "my thermometer",
   "name":"BertThermometer",
   "properties": {
      "organisation": "Geodan",
      "owner": "Bert"
   },
   "Locations@iot.navigationLink": "http://localhost:8080/v1/Things(5)/Locations",
   "Datastreams@iot.navigationLink": "http://localhost:8080/v1/Things(5)/Datastreams",
   "HistoricalLocations@iot.navigationLink": "http://localhost:8080/v1/Things(5)/HistoricalLocations"
}
```

## 2] Create a Location

Where is the physical or virtual object located?

Request: 

```sh
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "description": "my backyard",
    "name": "my backyard",
    "encodingType": "application/vnd.geo+json",
    "location": {
        "type": "Point",
        "coordinates": [-117.123,
        54.123]
    }
}' "http://localhost:8080/v1.0/Things(${Thing.@iot.id})/Locations"
```

Parameters: ${Thing.@iot.id} in the url path

Response:

```json
{
  "@iot.id": "1",
  "@iot.selfLink": "http://localhost:8080/v1.0/Locations(1)",
  "name": "my backyard",
  "description": "my backyard",
  "encodingType": "application/vnd.geo+json",
  "location": {
     "coordinates": [
        -117.123,
        54.123
     ],
     "type": "Point"
  },
  "Things": [
     {
        "@iot.id": "5"
     }
  ]
}
```

## 3] Create an ObservedProperty

What property is the sensor observing?

Request: 

```sh
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "air_temperature",
  "description": "Temperature of air in situ.",
  "definition": "http://mmisw.org/ont/ioos/parameter/air_temperature"
}' "http://localhost:8080/v1.0/ObservedProperties"
```

Parameters: -

Response:

```json
{
   "@iot.id": "2",
   "@iot.selfLink": "http://localhost:8080/v1.0/ObservedProperties(2)",
   "description": "Temperature of air in situ.",
   "name": "air_temperature",
   "definition": "http://mmisw.org/ont/ioos/parameter/air_temperature"
}
```

## 4] Create a Sensor

The sensor is what actually makes the observations. 

Request:

```sh
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{        
    "description": "Thermometer",
    "name": "thermometer",
    "encodingType": "application/pdf",
    "metadata": "https://en.wikipedia.org/wiki/Thermometer"
}' "http://localhost:8080/v1.0/Sensors"
```

Parameters: -

Response:

```json
{
   "@iot.id": "1",
   "@iot.selfLink": "http://localhost:8080/v1.0/Sensors(1)",
   "description": "Thermometer",
   "name": "thermometer",
   "encodingType": "application/pdf",
   "metadata": "https://en.wikipedia.org/wiki/Thermometer"
}
```

## 5] Create a Datastream

The Datastream represents a series of observations.

Request:

```sh
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "°C",
        "name": "degree Celsius",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Thermometer readings",
  "name": "Thermometer_readings",
  "Thing": {"@iot.id": ${Thing.@iot.id}},
  "ObservedProperty": {"@iot.id": ${ObservedProperty.@iot.id}},
  "Sensor": {"@iot.id": ${Sensor.@iot.id}}
}' "http://localhost:8080/v1.0/Datastreams"
```

Parameters: ${Sensor.@iot.id}, ${ObservedProperty.@iot.id}, ${Thing.@iot.id}

Response:

```json
{
   "@iot.id": "8",
   "@iot.selfLink": "http://localhost:8080/v1.0/Datastreams(8)",
   "description": "Thermometer readings",
   "name": "Thermometer_readings",
   "unitOfMeasurement": {
      "definition": "http://unitsofmeasure.org/ucum.html#para-30",
      "name": "degree Celsius",
      "symbol": "°C"
   },
   "observationType": "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement"
}
```

## 6] Add an Observation

Request: 

```sh
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "phenomenonTime": "2016-05-09T11:04:15.790Z",
  "resultTime" : "2016-05-09T11:04:15.790Z",
  "result" : 38,
  "Datastream":{"@iot.id":${Datastream.@iot.id}}
}' "http://localhost:8080/v1.0/Observations"
```

Parameters: ${Datastream.@iot.id}

Response:

```json
{
   "@iot.id": "2",
   "@iot.selfLink": "http://localhost:8080/v1.0/Observations(2)",
   "phenomenonTime": "2016-05-09T11:04:15.790Z",
   "result": 38,
   "resultTime": "2016-05-09T11:04:15.790Z",
   "Datastream@iot.navigationLink": "http://localhost:8080/v1.0/Observations(2)/Datastream",
   "FeatureOfInterest@iot.navigationLink": "http://localhost:8080/v1.0/Observations(185)/FeatureOfInterest"
}
```
