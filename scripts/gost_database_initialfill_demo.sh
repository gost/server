#!/usr/bin/env bash

if [[ $# -eq 0 ]] ; then
    echo "please supply the SensorThings host, for example: localhost:8080"
    exit 1
fi

HOST=$1

# CREATE THINGS AND LOCATIONS

# ----------------------------------------
# Create NodeMCU Aquarium Things(1) Locations(1)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "name": "NodeMCU Aquarium",
    "description": "NodeMCU DEVKIT V 0.9 Aquarium",
    "properties": {
		"id": "6BuE6ZSQ",
        "chip": "ESP2866",
        "owner": "Tim Ebben"
    },
	"Locations": [
		{
		    "name": "Geodan PK",
			"description": "2nd floor Geodan PK, inside Aquarium",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [4.913295, 52.343035]
			}
		}
	]
}' "$HOST/v1.0/Things"

# ----------------------------------------
# Create Netatmo PK Things(2) Locations(2)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "name": "Netatmo Weatherstation PK",
    "description": "Netatmo PK",
    "properties": {
        "outside_sn": "h035d52",
		"inside_mac": "70:ee:50:03:65:d4",
        "owner": "Tim Ebben"
    },
	"Locations": [
		{
		    "name": "Geodan PK",
			"description": "2nd floor Geodan PK, behind StevenB",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [4.913329, 52.343029]
			}
		}
	]
}' "$HOST/v1.0/Things"

# ----------------------------------------
# Create Netatmo VZ Things(3) Locations(3)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "name": "Netatmo Weatherstation VZ",
    "description": "Netatmo VZ",
    "properties": {
        "outside_sn": "h1e2998",
		"inside_mac": "70:ee:50:1d:f1:de",
        "owner": "KevinK"
    },
	"Locations": [
		{
		    "name": "Geodan VZ",
			"description": "2nd floor Geodan VZ",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [5.299545, 51.691778]
			}
		}
	]
}' "$HOST/v1.0/Things"

# ----------------------------------------
# Berts racing bike Things(4) Locations(4)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "name": "Racing bike Bert",
    "description": "Cycling sensor, used for demo purpose",
    "properties": {
        "owner": "Bert Temme"
    },
	"Locations": [
		{
		    "name": "Nobody knows",
			"description": "Its Everywhere",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [5.078448, 52.077259]
			}
		}
	]
}' "$HOST/v1.0/Things"

# ----------------------------------------
# SF BeeClear Things(5) Locations(5)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "name": "BeeClear SF",
    "description": "BeeClear connected to Stevens smart meter at home",
    "properties": {
        "owner": "StevenF"
    },
	"Locations": [
		{
		    "name": "Kombuis",
			"description": "StevenF house",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [4.8583665, 52.2777268]
			}
		}
	]
}' "$HOST/v1.0/Things"

# ----------------------------------------
# Zeecontainer Things(6) Locations(6)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "name": "Zeecontainer",
    "description": "Test zeecontainer beweging",
    "properties": {
        "owner": "EduardoD"
    },
	"Locations": [
		{
		    "name": "Sea",
			"description": "Cliffs of Dover &#9834; ",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [1.367776, 51.135059]
			}
		}
	]
}' "$HOST/v1.0/Things"

# ----------------------------------------
# Arduino The Things Network Things(7)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "name": "Arduino TTN",
    "description": "Arduino connected to TTN",
    "properties": {
        "owner": "StevenF"
    },
	"Locations": [
		{
		    "name": "Geodan PK",
			"description": "2nd floor Geodan PK, inside Aquarium",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [4.913295, 52.343035]
			}
		}
	]
}' "$HOST/v1.0/Things"

# ----------------------------------------
# Alexanders NodeMCU connected to various sensors Things(8) 
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "name": "Alexanders NodeMCU connected to various sensors",
    "description": "Alexanders NodeMCU connected to various sensors",
    "properties": {
        "owner": "AlexanderF"
    },
	"Locations": [
		{
		    "name": "Geodan PK",
			"description": "2nd floor Geodan PK, inside Aquarium",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [4.913295, 52.343035]
			}
		}
	]
}' "$HOST/v1.0/Things"

# ----------------------------------------
# Eduardo's phone Things(9) 
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json"  -d '{
    "name": "EduardoD phone",
    "description": "EduardoD phone",
    "properties": {
        "owner": "EduardoD"
    },
	"Locations": [
		{
		    "name": "Geodan PK",
			"description": "2nd floor Geodan PK, inside Aquarium",
			"encodingType": "application/vnd.geo+json",
			"location": {
				"type": "Point",
				"coordinates": [4.913295, 52.343035]
			}
		}
	]
}' "$HOST/v1.0/Things"


# CREATE OBSERVED PROPERTIES


# ----------------------------------------
# Temperature ObservedProperties(1)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "air_temperature",
  "description": "Temperature of air in situ.",
  "definition": "http://mmisw.org/ont/ioos/parameter/air_temperature"
}' "$HOST/v1.0/ObservedProperties"

# ----------------------------------------
# Humidity ObservedProperties(2)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "air_humidity",
  "description": "Humidity of air in situ.",
  "definition": "http://mmisw.org/ont/ioos/parameter/humidity"
}' "$HOST/v1.0/ObservedProperties"

# ----------------------------------------
# Air pressure ObservedProperties(3)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "air_pressure",
  "description": "Pressure exerted by overlying air",
  "definition": "http://mmisw.org/ont/ioos/parameter/air_pressure"
}' "$HOST/v1.0/ObservedProperties"

# ----------------------------------------
# CO2 ObservedProperties(4)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "CO2",
  "description": "CO2 in situ",
  "definition": "http://mmisw.org/ont/ioos/parameter/co2"
}' "$HOST/v1.0/ObservedProperties"

# ----------------------------------------
# Sound ObservedProperties(5)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "sound",
  "description": "Sound in situ",
  "definition": "http://mmisw.org/ont/ioos/parameter/sound"
}' "$HOST/v1.0/ObservedProperties"

# ----------------------------------------
# Speed ObservedProperties(6)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "speed",
  "description": "movement speed",
  "definition": "http://mmisw.org/ont/ioos/parameter/speed"
}' "$HOST/v1.0/ObservedProperties"

# ----------------------------------------
# Electricity ObservedProperties(7)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "electricity",
  "description": "electricity usage",
  "definition": "http://mmisw.org/ont/ioos/parameter/electricity"
}' "$HOST/v1.0/ObservedProperties"

# ----------------------------------------
# Acceleration ObservedProperties(8)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "Acceleration",
  "description": "Acceleration",
  "definition": "http://www.qudt.org/qudt/owl/1.0.0/quantity/Instances.html"
}' "$HOST/v1.0/ObservedProperties"

# ----------------------------------------
# Rotation ObservedProperties(9)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "name": "Rotation",
  "description": "Rotation in degrees along the different axis",
  "definition": "http://sensorml.com/ont/swe/property/SensorOrientation"
}' "$HOST/v1.0/ObservedProperties"

# CREATE SENSORS

# ----------------------------------------
# HTU21D Sensors(1)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "name": "HTU21D",
    "description": "HTU21D for sensing temperature and humidity",
    "encodingType": "application/pdf",
    "metadata": "https://cdn-shop.adafruit.com/datasheets/1899_HTU21D.pdf"
}' "$HOST/v1.0/Sensors"

# ----------------------------------------
# Netatmo Sensors(2)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "name": "Netatmo",
    "description": "Netatmo Weatherstation",
    "encodingType": "application/pdf",
    "metadata": "https://www.netatmo.com/product/station"
}' "$HOST/v1.0/Sensors"

# ----------------------------------------
# Garmin speed sensor Sensors(3)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "name": "Garmin Speed Sensor",
    "description": "Cadanssensor",
    "encodingType": "application/pdf",
    "metadata": "https://www8.garmin.com/manuals/webhelp/edge520/EN-US/GUID-E7E492A5-5342-4B27-875B-0C7B5D5F14E7.html"
}' "$HOST/v1.0/Sensors"

# ----------------------------------------
# Slimme meter P1 (BeeClear) Sensors(4)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "name": "BeeClear",
    "description": "BeeClear Energiemanager V2",
    "encodingType": "application/pdf",
    "metadata": "https://beeclear.nl/docs/handleiding.pdf"
}' "$HOST/v1.0/Sensors"

# ----------------------------------------
# Accelerometer Sensors(5)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "name": "Accelerometer",
    "description": "Accelerometer",
    "encodingType": "application/pdf",
    "metadata": "Calibration date: Jan 1, 2014"
}' "$HOST/v1.0/Sensors"

# ----------------------------------------
# Gyroscope Sensors(6)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "name": "Gyroscope",
    "description": "Gyroscope",
    "encodingType": "application/pdf",
    "metadata": "Calibration date: Jan 1, 2014"
}' "$HOST/v1.0/Sensors"

# ----------------------------------------
# MQ135 Sensors(7)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "name": "MQ135",
    "description": "MQ135 Ait Quality Sensor. Sensitive for Benzene, Alcohol, smoke can be used to calculate CO2",
    "encodingType": "application/pdf",
    "metadata": "https://www.olimex.com/Products/Components/Sensors/SNS-MQ135/resources/SNS-MQ135.pdf"
}' "$HOST/v1.0/Sensors"

# ----------------------------------------
# DHT22 Sensors(8)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "name": "DHT22",
    "description": "DHT22 for sensing temperature and humidity",
    "encodingType": "application/pdf",
    "metadata": "https://www.sparkfun.com/datasheets/Sensors/Temperature/DHT22.pdf"
}' "$HOST/v1.0/Sensors"

# ----------------------------------------
# BMP180 Sensors(9)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "name": "BMP180",
    "description": "BMP180 for sensing barometric pressure/temperature/altitude",
    "encodingType": "application/pdf",
    "metadata": "https://cdn-shop.adafruit.com/datasheets/BST-BMP180-DS000-09.pdf"
}' "$HOST/v1.0/Sensors"

# ----------------------------------------
# MPU-6050 Sensors(10)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "name": "MPU-6050",
    "description": "MPU-6050 3 axis gyro & 3 axis accelerometer",
    "encodingType": "application/pdf",
    "metadata": "https://www.cdiweb.com/datasheets/invensense/MPU-6050_DataSheet_V3%204.pdf"
}' "$HOST/v1.0/Sensors"

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
  "name": "NodeMCU Temperature PK Aquarium",
  "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
  "description": "NodeMCU Temperature readings indoor",
  "Thing": {"@iot.id": 1},
  "ObservedProperty": {"@iot.id": 1},
  "Sensor": {"@iot.id": 1}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# NodeMCU Humidity Datastreams(2)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "%",
        "name": "humidity",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "NodeMCU Humidity PK Aquarium",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "NodeMCU humidity readings indoor",
    "Thing": {"@iot.id": 1},
    "ObservedProperty": {"@iot.id": 2},
    "Sensor": {"@iot.id": 1}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Temperature indoor Datastreams(3)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "°C",
        "name": "degree Celsius",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo PK temperature 2nd floor",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo PK temperature readings indoor",
    "Thing": {"@iot.id": 2},
    "ObservedProperty": {"@iot.id": 1},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Humidity indoor Datastreams(4)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "%",
        "name": "humidity",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo PK humidity 2nd floor",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo PK humidity readings indoor",
    "Thing": {"@iot.id": 2},
    "ObservedProperty": {"@iot.id": 2},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Air pressure indoor Datastreams(5)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "mb",
        "name": "bar",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo PK pressure 2nd floor",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo PK pressure readings indoor",
    "Thing": {"@iot.id": 2},
    "ObservedProperty": {"@iot.id": 3},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo CO2 indoor Datastreams(6)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "ppm",
        "name": "parts per million",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo PK CO2 2nd floor",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo PK CO2 readings indoor",
    "Thing": {"@iot.id": 2},
    "ObservedProperty": {"@iot.id": 4},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Sound indoor Datastreams(7)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "db",
        "name": "decibel",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo PK sound 2nd floor",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo sound readings indoor",
    "Thing": {"@iot.id": 2},
    "ObservedProperty": {"@iot.id": 5},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Temperature outdoor Datastreams(8)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "°C",
        "name": "degree Celsius",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo PK temperature outside",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo PK temperature readings outdoor",
    "Thing": {"@iot.id": 2},
    "ObservedProperty": {"@iot.id": 1},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Humidity outdoor Datastreams(9)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "%",
        "name": "humidity",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo PK humidity outside",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo PK humidity readings outdoor",
    "Thing": {"@iot.id": 2},
    "ObservedProperty": {"@iot.id": 2},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Temperature indoor Datastreams(10)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "°C",
        "name": "degree Celsius",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo VZ temperature 2nd floor",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo VZ temperature readings indoor",
    "Thing": {"@iot.id": 3},
    "ObservedProperty": {"@iot.id": 1},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Humidity indoor Datastreams(11)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "%",
        "name": "humidity",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo VZ humidity 2nd floor",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo VZ humidity readings indoor",
    "Thing": {"@iot.id": 3},
    "ObservedProperty": {"@iot.id": 2},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Air pressure indoor Datastreams(12)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "mb",
        "name": "bar",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo VZ pressure 2nd floor",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo VZ pressure readings indoor",
    "Thing": {"@iot.id": 3},
    "ObservedProperty": {"@iot.id": 3},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo CO2 indoor Datastreams(13)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "ppm",
        "name": "parts per million",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo VZ CO2 2nd floor",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo VZ CO2 readings indoor",
    "Thing": {"@iot.id": 3},
    "ObservedProperty": {"@iot.id": 4},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Sound indoor Datastreams(14)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "db",
        "name": "decibel",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo VZ sound 2nd floor",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo VZ sound readings indoor",
    "Thing": {"@iot.id": 3},
    "ObservedProperty": {"@iot.id": 5},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Temperature outdoor Datastreams(15)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "°C",
        "name": "degree Celsius",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo VZ temperature outside",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo VZ temperature readings outdoor",
    "Thing": {"@iot.id": 3},
    "ObservedProperty": {"@iot.id": 1},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Netatmo Humidity outdoor Datastreams(16)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "%",
        "name": "humidity",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Netatmo VZ humidity outside",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Netatmo VZ humidity readings outdoor",
    "Thing": {"@iot.id": 3},
    "ObservedProperty": {"@iot.id": 2},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Speed measurements Datastreams(17)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "km/h",
        "name": "kilometers per hour",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Berts racing bike measurements",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Speed measurements Berts racing bike",
    "Thing": {"@iot.id": 4},
    "ObservedProperty": {"@iot.id": 6},
    "Sensor": {"@iot.id": 3}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Energy usage measurements Datastreams(18)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "kW",
        "name": "Kilo Watt",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Energy usage SF",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Energy usage measurements SF from BeeClear",
    "Thing": {"@iot.id": 5},
    "ObservedProperty": {"@iot.id": 7},
    "Sensor": {"@iot.id": 4}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Container acceleration x-axis measurements Datastreams(19)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
		"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html",
		"name": "g",
		"symbol": "g"
	},
    "name": "Container acceleration x-axis",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Container acceleration x-axis",
    "Thing": {"@iot.id": 6},
    "ObservedProperty": {"@iot.id": 8},
    "Sensor": {"@iot.id": 5}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Container acceleration y-axis measurements Datastreams(20)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
		"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html",
		"name": "g",
		"symbol": "g"
	},
    "name": "Container acceleration y-axis",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Container acceleration y-axis",
    "Thing": {"@iot.id": 6},
    "ObservedProperty": {"@iot.id": 8},
    "Sensor": {"@iot.id": 5}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Container acceleration z-axis measurements Datastreams(21)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
		"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html",
		"name": "g",
		"symbol": "g"
	},
    "name": "Container acceleration z-axis",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Container acceleration z-axis",
    "Thing": {"@iot.id": 6},
    "ObservedProperty": {"@iot.id": 9},
    "Sensor": {"@iot.id": 6}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Container rotation x-axis measurements Datastreams(22)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
		"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html",
		"name": "g",
		"symbol": "g"
	},
    "name": "Container rotation x-axis",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Container rotation x-axis",
    "Thing": {"@iot.id": 6},
    "ObservedProperty": {"@iot.id": 9},
    "Sensor": {"@iot.id": 6}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Container rotation y-axis measurements Datastreams(23)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
		"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html",
		"name": "g",
		"symbol": "g"
	},
    "name": "Container rotation y-axis",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Container rotation y-axis",
    "Thing": {"@iot.id": 6},
    "ObservedProperty": {"@iot.id": 9},
    "Sensor": {"@iot.id": 6}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Container rotation z-axis measurements Datastreams(24)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
		"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html",
		"name": "g",
		"symbol": "g"
	},
    "name": "Container rotation z-axis",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Container rotation z-axis",
    "Thing": {"@iot.id": 6},
    "ObservedProperty": {"@iot.id": 8},
    "Sensor": {"@iot.id": 5}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Arduino TTN Temperature (25)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
        "symbol": "°C",
        "name": "degree Celsius",
        "definition": "http://unitsofmeasure.org/ucum.html#para-30"
    },
    "name": "Arduino TTN temperature",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Arduino TTN temperature measurements",
    "Thing": {"@iot.id": 7},
    "ObservedProperty": {"@iot.id": 1},
    "Sensor": {"@iot.id": 2}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Arduino TTN CO2 (26)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
		"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html",
		"name": "Parts Per Million",
		"symbol": "PPM"
	},
    "name": "Arduino TTN CO2",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Arduino TTN CO2 measurements",
    "Thing": {"@iot.id": 7},
    "ObservedProperty": {"@iot.id": 4},
    "Sensor": {"@iot.id": 7}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Alexander Orientation of the MPU-6050 (27)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
		"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html",
		"name": "Degrees",
		"symbol": "Degree"
	},
    "name": "Orientation of the MPU-6050",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "Orientation of the MPU-6050",
    "Thing": {"@iot.id": 8},
    "ObservedProperty": {"@iot.id": 9},
    "Sensor": {"@iot.id": 10}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Eduardo phone rotation x-axis measurements Datastreams(28)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
		"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html",
		"name": "degrees",
		"symbol": "Deg"
	},
    "name": "EduardoD phone rotation x-axis",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "EduardoD phone rotation x-axis",
    "Thing": {"@iot.id": 9},
    "ObservedProperty": {"@iot.id": 9},
    "Sensor": {"@iot.id": 6}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Eduardo phone rotation y-axis measurements Datastreams(29)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
		"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html",
		"name": "degrees",
		"symbol": "Deg"
	},
    "name": "EduardoD phone rotation y-axis",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "EduardoD phone rotation y-axis",
    "Thing": {"@iot.id": 9},
    "ObservedProperty": {"@iot.id": 9},
    "Sensor": {"@iot.id": 6}
}' "$HOST/v1.0/Datastreams"

# ----------------------------------------
# Eduardo phone rotation z-axis measurements Datastreams(30)
# ----------------------------------------
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
    "unitOfMeasurement": {
		"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html",
		"name": "degrees",
		"symbol": "Deg"
	},
    "name": "EduardoD phone rotation z-axis",
    "observationType":"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement",
    "description": "EduardoD phone rotation z-axis",
    "Thing": {"@iot.id": 9},
    "ObservedProperty": {"@iot.id": 9},
    "Sensor": {"@iot.id": 6}
}' "$HOST/v1.0/Datastreams"

