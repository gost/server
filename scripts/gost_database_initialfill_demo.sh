#!/usr/bin/env bash

# CREATE THINGS AND LOCATIONS

# ----------------------------------------
# Create NodeMCU Things(1) Locations(1)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "description": "NodeMCU DEVKIT 0.9",
    "properties": {
		"id": "6BuE6ZSQ",
        "chip": "ESP2866",
        "owner": "Tim Ebben"
    },
	"Locations": [
		{
			"description": "Geodan PK",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [4.913021, 52.342417]
			}
		}
	]
}' "http://gost.geodan.nl/v1.0/Things"

# ----------------------------------------
# Create Netatmo Things(2) Locations(2)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "description": "Netatmo Weatherstation VZ",
    "properties": {
        "outside_sn": "h035d52",
		"inside_mac": "70:ee:50:03:65:d4",
        "owner": "Tim Ebben"
    },
	"Locations": [
		{
			"description": "Geodan VZ - 2nd",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [5.299574, 51.691786]
			}
		}
	]
}' "http://gost.geodan.nl/v1.0/Things"

# ----------------------------------------
# Berts racing bike Things(3) Locations(3)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "description": "Racing bike Bert",
    "properties": {
        "owner": "Bert Temme"
    },
	"Locations": [
		{
			 "description": "Its Everywhere",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [5.078448, 52.077259]
			}
		}
	]
}' "http://gost.geodan.nl/v1.0/Things"

# ----------------------------------------
# SF BeeClear Things(4) Locations(4)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "description": "BeeClear SF",
    "properties": {
        "owner": "Steven"
    },
	"Locations": [
		{
			"description": "Kombuis",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [4.8583665, 52.2777268]
			}
		}
	]
}' "http://gost.geodan.nl/v1.0/Things"


# CREATE OBSERVED PROPERTIES


# ----------------------------------------
# Temperature ObservedProperties(1)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "air_temperature",
  "description": "Temperature of air in situ.",
  "definition": "http://mmisw.org/ont/ioos/parameter/air_temperature"
}' "http://gost.geodan.nl/v1.0/ObservedProperties"

# ----------------------------------------
# Humidity ObservedProperties(2)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "air_humidity",
  "description": "Humidity of air in situ.",
  "definition": "http://mmisw.org/ont/ioos/parameter/humidity"
}' "http://gost.geodan.nl/v1.0/ObservedProperties"

# ----------------------------------------
# Air pressure ObservedProperties(3)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "air_pressure",
  "description": "Pressure exerted by overlying air",
  "definition": "http://mmisw.org/ont/ioos/parameter/air_pressure"
}' "http://gost.geodan.nl/v1.0/ObservedProperties"

# ----------------------------------------
# CO2 ObservedProperties(4)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "CO2",
  "description": "CO2 in situ",
  "definition": "http://mmisw.org/ont/ioos/parameter/co2"
}' "http://gost.geodan.nl/v1.0/ObservedProperties"

# ----------------------------------------
# Sound ObservedProperties(5)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "sound",
  "description": "Sound in situ",
  "definition": "http://mmisw.org/ont/ioos/parameter/sound"
}' "http://gost.geodan.nl/v1.0/ObservedProperties"

# ----------------------------------------
# Speed ObservedProperties(6)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "speed",
  "description": "movement speed",
  "definition": "http://mmisw.org/ont/ioos/parameter/speed"
}' "http://gost.geodan.nl/v1.0/ObservedProperties"

# ----------------------------------------
# Electricity ObservedProperties(7)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "electricity",
  "description": "electricity usage",
  "definition": "http://mmisw.org/ont/ioos/parameter/electricity"
}' "http://gost.geodan.nl/v1.0/ObservedProperties"


# CREATE SENSORS


# ----------------------------------------
# HTU21D Sensors(1)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{        
    "description": "HTU21D",
    "encodingType": "application/pdf",
    "metadata": "https://cdn-shop.adafruit.com/datasheets/1899_HTU21D.pdf"
}' "http://gost.geodan.nl/v1.0/Sensors"

# ----------------------------------------
# Netatmo Sensors(2)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{        
    "description": "Netatmo",
    "encodingType": "application/pdf",
    "metadata": "https://www.netatmo.com/product/station"
}' "http://gost.geodan.nl/v1.0/Sensors"

# ----------------------------------------
# Garmin speed sensor Sensors(3)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{        
    "description": "Garmin Speed Sensor",
    "encodingType": "application/pdf",
    "metadata": "https://www8.garmin.com/manuals/webhelp/edge520/EN-US/GUID-E7E492A5-5342-4B27-875B-0C7B5D5F14E7.html"
}' "http://gost.geodan.nl/v1.0/Sensors"

# ----------------------------------------
# Slimme meter P1 (BeeClear) Sensors(4)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{        
    "description": "BeeClear Energiemanager V2",
    "encodingType": "application/pdf",
    "metadata": "https://beeclear.nl/docs/handleiding.pdf"
}' "http://gost.geodan.nl/v1.0/Sensors"


# CREATE DATASTREAMS


# ----------------------------------------
# NodeMCU Temperature Datastreams(1)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "°C",
        "name": "degree Celsius",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "NodeMCU Temperature readings indoor",
  "Thing": {"@iot.id": 1},
  "ObservedProperty": {"@iot.id": 1},
  "Sensor": {"@iot.id": 1}
}' "http://gost.geodan.nl/v1.0/Datastreams"

# ----------------------------------------
# NodeMCU Humidity Datastreams(2)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "%",
        "name": "humidity",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "NodeMCU humidity readings indoor",
  "Thing": {"@iot.id": 1},
  "ObservedProperty": {"@iot.id": 2},
  "Sensor": {"@iot.id": 1}
}' "http://gost.geodan.nl/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Temperature indoor Datastreams(3)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "°C",
        "name": "degree Celsius",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Netatmo temperature readings indoor",
  "Thing": {"@iot.id": 2},
  "ObservedProperty": {"@iot.id": 1},
  "Sensor": {"@iot.id": 2}
}' "http://gost.geodan.nl/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Humidity indoor Datastreams(4)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "%",
        "name": "humidity",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Netatmo humidity readings indoor",
  "Thing": {"@iot.id": 2},
  "ObservedProperty": {"@iot.id": 2},
  "Sensor": {"@iot.id": 2}
}' "http://gost.geodan.nl/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Air pressure indoor Datastreams(5)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "mb",
        "name": "bar",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Netatmo pressure readings indoor",
  "Thing": {"@iot.id": 2},
  "ObservedProperty": {"@iot.id": 3},
  "Sensor": {"@iot.id": 2}
}' "http://gost.geodan.nl/v1.0/Datastreams"

# ----------------------------------------
# Netatmo CO2 indoor Datastreams(6)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "ppm",
        "name": "parts per million",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Netatmo CO2 readings indoor",
  "Thing": {"@iot.id": 2},
  "ObservedProperty": {"@iot.id": 4},
  "Sensor": {"@iot.id": 2}
}' "http://gost.geodan.nl/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Sound indoor Datastreams(7)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "db",
        "name": "decibel",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Netatmo sound readings indoor",
  "Thing": {"@iot.id": 2},
  "ObservedProperty": {"@iot.id": 5},
  "Sensor": {"@iot.id": 2}
}' "http://gost.geodan.nl/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Temperature outdoor Datastreams(8)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "°C",
        "name": "degree Celsius",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Netatmo temperature readings outdoor",
  "Thing": {"@iot.id": 2},
  "ObservedProperty": {"@iot.id": 1},
  "Sensor": {"@iot.id": 2}
}' "http://gost.geodan.nl/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Humidity outdoor Datastreams(9)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "%",
        "name": "humidity",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Netatmo humidity readings outdoor",
  "Thing": {"@iot.id": 2},
  "ObservedProperty": {"@iot.id": 2},
  "Sensor": {"@iot.id": 2}
}' "http://gost.geodan.nl/v1.0/Datastreams"

# ----------------------------------------
# Speed measurements Datastreams(10)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "km/h",
        "name": "kilometers per hour",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Speed measurements Berts racing bike",
  "Thing": {"@iot.id": 3},
  "ObservedProperty": {"@iot.id": 6},
  "Sensor": {"@iot.id": 3}
}' "http://gost.geodan.nl/v1.0/Datastreams"

# ----------------------------------------
# Energy usage measurements Datastreams(11)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "kW",
        "name": "Kilo Watt",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "Energy usage measurements SF",
  "Thing": {"@iot.id": 4},
  "ObservedProperty": {"@iot.id": 7},
  "Sensor": {"@iot.id": 4}
}' "http://gost.geodan.nl/v1.0/Datastreams"

