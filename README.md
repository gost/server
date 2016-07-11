<img src="src/client/assets/img/icon.png" width="353"><br />
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Geodan/gost/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/Geodan/gost?status.svg)](https://godoc.org/github.com/Geodan/gost)
[![Build Status](http://beta.drone.io/api/badges/drone/drone/status.svg)](https://drone.io/github.com/Geodan/gost/latest)
[![Go Report Card](https://goreportcard.com/badge/geodan/gost)](https://goreportcard.com/report/geodan/gost)
[![Coverage Status](https://coveralls.io/repos/github/Geodan/gost/badge.svg?branch=master)](https://coveralls.io/github/Geodan/gost?branch=master)
[![Join the chat at https://gitter.im/Geodan/gost](https://badges.gitter.im/Geodan/gost.svg)](https://gitter.im/Geodan/gost?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)<br />

GOST (Go-SensorThings) is an IoT Platform written in Golang (Go). It implements the Sensing profile (part 1) of the [OGC SensorThings API] (http://ogc-iot.github.io/ogc-iot-api/api.html) standard including the MQTT extension.

Implementation of the Tasking profile (part 2) and Rules Engine profile (part 3) of the OGC SensorThings API is planned as a future work activity.

## Disclaimer

GOST is alpha software and is not (yet) considered appropriate for customer use. Feel free to help development :-)

## Binaries

Release 0.2: 2016-06-10

[https://github.com/Geodan/gost/releases/tag/0.2] (https://github.com/Geodan/gost/releases/tag/0.2)

Binaries are build for Windows, Ubuntu and OSX.

## News

2016-05-30: Testing started against the OGC SensorThings API Test Suite 

## OGC Compliance testing

Gost is being tested against the OGC SensorThings API Test Suite [https://github.com/opengeospatial/ets-sta10](https://github.com/opengeospatial/ets-sta10)

Gost Compliance Testing Status:

| Conformance Class                     | Reference | Implementation status |Test Status               |
|---------------------------------------|-----------|-----------------------|--------------------------| 
| Sensing Core                          | A.1       | alpha                 | 6 passed, 0 failed       |
| Filtering Extension                   | A.2       | alpha                 | Testing not started      |
| Create-Update-Delete                  | A.3       | alpha                 | 8 passed, 1 failed       |
| Batch Request                         | A.4       | -                     | Tests not implemented    |
| Sensing MultiDatastream Extension     | A.5       | -                     | Tests not implemented    |
| Sensing Data Array Extension          | A.6       | -                     | Tests not implemented    |
| MQTT Extension for Create and Update  | A.7       | alpha                 | Tests not implemented    |
| MQTT Extension for Receiving Updates  | A.8       | alpha                 | Tests not implemented    |

## Installation and configuration

[Gost installation](docs/gost_installation.md)

[Gost configuration](docs/gost_configuration.md)

## Samples

HTTP Api: For sample requests (setting up sensors/datastreams/things and adding observations) see the tests in the [playground](test/playground_tests.md). 
For a complete collection of working requests install Postman and import the [Postman file](test/GOST.json.postman_collection) 

MQTT: For getting started with Gost and MQTT for publishing/receiving data see [Gost and MQTT - Getting started](docs/gost_mqtt_getting_started.md)

## Dependencies

[yaml v2](https://github.com/go-yaml/yaml)<br />
[pq](https://github.com/lib/pq)<br />
[mux](https://github.com/gorilla/mux)<br />
[Paho](https://github.com/eclipse/paho.mqtt.golang)<br />

## Roadmap

- Complete implementation of the OGC SensorThings spec
- Test coverage!
- Frontend
- Benchmarks
- Authentication
- Different storage providers such as MongoDB (Now using PostgreSQL)

## TODO

[See wiki](https://github.com/Geodan/gost/wiki/TODO)

## Benchmarks

[Publish observation (MQTT)](https://github.com/Geodan/gost/wiki/Benchmark---publish-observation-(MQTT))
