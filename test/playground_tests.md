## Playground tests

Here a collection of sample requests to start working with the GOST API. It contains 6 steps: create a Thing, create a Location, create an ObservedProperty, create a Sensor, create a Datastream and create an Observation.

Instructions:

. replace http://localhost:8080 with the url of your own test server

. when something is created with a HTTP Post request, and Id is returned. Uses this Id in subsequent requests (see parameters)

## 1] Create a Thing

A thing is an object in either the physical or virtual world.

Request:

```sh
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "description": "my thermometer",
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
   "properties": {
      "organisation": "Geodan",
      "owner": "Bert"
   },
   "Locations@iot.navigationLink": "../Things(5)/Locations",
   "Datastreams@iot.navigationLink": "../Things(5)/Datastreams",
   "HistoricalLocations@iot.navigationLink": "../Things(5)/HistoricalLocations"
}
```

## 2] Create a Location

Where is the physical or virtual object located?

Request: 

```sh
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "description": "my backyard",
    "encodingType": "application/vnd.geo+json",
    "location": {
        "type": "Point",
        "coordinates": [-117.123,
        54.123]
    }
}' "http://localhost:8080/v1.0/Locations"
```

Parameters: Things.@iot.id

Response:

```json
{
  "@iot.id": "1",
  "description": "my backyard",
  "encodingtype": "application/vnd.geo+json",
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
    "encodingType": "application/pdf",
    "metadata": "https://en.wikipedia.org/wiki/Thermometer"
}' "http://localhost:8080/v1.0/Sensors"
```

Parameters: -

Response:

```json
{
   "@iot.id": "1",
   "description": "Thermometer",
   "encodingtype": "application/pdf",
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
        "definition": "http://unitsofmeasure.org/ucum.html#para-30
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Thermometer readings",
  "Thing": {"@iot.id": "1"},
  "ObservedProperty": {"@iot.id": "1"},
  "Sensor": {"@iot.id": "1"}
}' "http://localhost:8080/v1.0/Datastreams"
```

Parameters: Sensor.@iot.id, ObservedProperty.@iot.id, Thing.@iot.id

Response:

```json
{
   "@iot.id": "8",
   "description": "Thermometer readings",
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
curl -X POST -d " { \"phenomenonTime\": \"2016-05-09T11:04:15.790Z\", \"result\": 20 }" --header "Content-Type:application/json" http://localhost:8080/v1.0/Datastreams(1)/Observations
```

Parameters:
Datastream.@iot.id in the url path

Response:

```json
{
   "@iot.id": "2",
   "@iot.selfLink": "http://localhost:8080/v1.0/Observations(2)",
   "phenomenonTime": "2016-05-09T11:04:15.790Z",
   "result": 20,
   "resultTime": "NULL",
   "Datastream@iot.navigationLink": "../Observations(2)/Datastream",
   "FeatureOfInterest@iot.navigationLink": "../Observations(2)/FeatureOfInterest"
}
```
