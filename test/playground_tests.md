## Playground tests

Here a collection of sample requests to start working with the GOST API. It contains 6 steps: create a Thing, create a Location, create an ObservedProperty, create a Sensor, create a Datastream and create an Observation.

## 1] Create a Thing

Request:

```sh
curl -X POST -d "{\"description\": \"my thermometer\"}" --header "Content-Type:application/json" http://localhost:8080/v1.0/Things
```

Parameters: -

Response:
```json
{
   "@iot.id": "5",
   "@iot.selfLink": "http://localhost:8080/v1.0/Things(5)",
   "description": "my thermometer",
   "Locations@iot.navigationLink": "../Things(5)/Locations",
   "Datastreams@iot.navigationLink": "../Things(5)/Datastreams",
   "HistoricalLocations@iot.navigationLink": "../Things(5)/HistoricalLocations"
}
```

## 2] Create a Location

Request: 

```sh
curl -X POST -d "{\"description\":\"My house\",\"encodingType\":\"application/vnd.geo+json\",\"location\":{\"type\":\"Point\",\"coordinates\":[-114.13046,51.083694]},\"Things\":[{\"@iot.id\":5}]}" --header "Content-Type:application/json" http://localhost:8080/v1.0/Locations
```

Parameters: Things.@iot.id

Response:

```json
{
  "@iot.id": "1",
  "description": "My house",
  "encodingtype": "application/vnd.geo+json",
  "location": {
     "coordinates": [
        -114.13046,
        51.083694
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

Request: 

```sh
curl -X POST -d "{\"name\": \"air_temperature\", \"definition\": \"http://mmisw.org/ont/ioos/parameter/air_temperature\", \"description\": \"Temperature of air in situ.\"}" --header "Content-Type:application/json" http://localhost:8080/v1.0/ObservedProperties
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

Request:

```sh
curl -X POST -d "{ \"description\": \"Thermometer\", \"encodingType\": \"text/html\", \"metadata\": \"https://en.wikipedia.org/wiki/Thermometer\" }" --header "Content-Type:application/json" http://localhost:8080/v1.0/Sensors
```

Parameters: -

Response:

```json
{
   "@iot.id": "1",
   "description": "Thermometer",
   "encodingtype": "text/html",
   "metadata": "https://en.wikipedia.org/wiki/Thermometer"
}
```

## 5] Create a Datastream

Request:

```sh
curl -X POST -d "{ \"description\": \"Thermometer readings\", \"observationType\": \"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement\", \"unitOfMeasurement\": { \"name\": \"degree Celsius\", \"symbol\": \"°C\", \"definition\": \"http://unitsofmeasure.org/ucum.html#para-30\" }, \"Sensor\": { \"@iot.id\": \"1\" }, \"ObservedProperty\": { \"@iot.id\": \"1\" }, \"Thing\": { \"@iot.id\": \"5\" } }" --header "Content-Type:application/json" http://localhost:8080/v1.0/Datastreams
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
      "symbol": "∩┐╜C"
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
