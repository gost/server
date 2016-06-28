curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "description": "my thermometer",
    "properties": {
        "organisation": "My Organisation",
        "owner": "Me"
    }
}' "http://gost.geodan.nl/v1.0/Things"

curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "description": "my backyard",
    "encodingType": "application/vnd.geo+json",
    "location": {
        "type": "Point",
        "coordinates": [-117.123,
        54.123]
    }
}' "http://gost.geodan.nl/v1.0/Things(1)/Locations"

curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "air_temperature",
  "description": "Temperature of air in situ.",
  "definition": "http://mmisw.org/ont/ioos/parameter/air_temperature"
}' "http://gost.geodan.nl/v1.0/ObservedProperties"

curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{        
    "description": "Thermometer",
    "encodingType": "application/pdf",
    "metadata": "https://en.wikipedia.org/wiki/Thermometer"
}' "http://gost.geodan.nl/v1.0/Sensors"

curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "°C",
        "name": "degree Celsius",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Thermometer readings",
  "Thing": {"@iot.id": 1},
  "ObservedProperty": {"@iot.id": 1},
  "Sensor": {"@iot.id": 1}
}' "http://gost.geodan.nl/v1.0/Datastreams"


curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "phenomenonTime": "2016-05-09T11:04:15.790Z",
  "resultTime" : "2016-05-09T11:04:15.790Z",
  "result" : 38,
  "Datastream":{"@iot.id":1}
}' "http://gost.geodan.nl/v1.0/Observations"


curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "description": "Underground Air Quality in NYC train tunnels",
    "encodingType": "application/vnd.geo+json",
    "feature": {
        "coordinates": [51.08386,-114.13036],
        "type": "Point"
      }
}' "http://gost.geodan.nl/v1.0/FeaturesOfInterest"
