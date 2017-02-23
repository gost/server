<img src="src/client/assets/img/icon.png" width="353"><br />
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Geodan/gost/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/Geodan/gost?status.svg)](https://godoc.org/github.com/Geodan/gost)
[![Build Status](http://beta.drone.io/api/badges/drone/drone/status.svg)](https://drone.io/github.com/Geodan/gost/latest)
[![Go Report Card](https://goreportcard.com/badge/geodan/gost)](https://goreportcard.com/report/geodan/gost)
[![Coverage Status](https://coveralls.io/repos/github/Geodan/gost/badge.svg?branch=master)](https://coveralls.io/github/Geodan/gost?branch=master)
[![Join the chat at https://gitter.im/Geodan/gost](https://badges.gitter.im/Geodan/gost.svg)](https://gitter.im/Geodan/gost?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)<br />

GOST (Go-SensorThings) is an IoT Platform written in Golang (Go). It implements the Sensing profile (part 1) of the [OGC SensorThings API] (http://ogc-iot.github.io/ogc-iot-api/api.html) standard including the MQTT extension.

(http://www.opengeospatial.org/resource/products/details/?pid=1419)[https://portal.opengeospatial.org/public_ogc/compliance/srv_ogc_compliance_badge.php?id=118&pid=1419]

Implementation of the Tasking profile (part 2) and Rules Engine profile (part 3) of the OGC SensorThings API is planned as a future work activity.

The GOST website and blog can be found at [www.gostserver.xyz](https://www.gostserver.xyz).

## Disclaimer

GOST is alpha software and is not (yet) considered appropriate for customer use. Feel free to help development :-)

## Binaries

Current release 0.3: 2016-12-06

[https://github.com/Geodan/gost/releases/tag/0.3] (https://github.com/Geodan/gost/releases/tag/0.3)

Binaries are build for Windows, Ubuntu and OSX.


## Roadmap

| Date       	|             Version 	| Features                                                        	|
|------------	|---------------------	|-----------------------------------------------------------------	|
| 2017-02-06 	| 0.4                 	| OGC Test level 3 compliant                                      	|
| 2017-03-06 	| 0.5                 	| TBD                                                             	|


## Docker support

See [GOST and Docker](docs/gost_docker.md)

TL;DR:

```
$ wget https://raw.githubusercontent.com/Geodan/gost/master/src/docker-compose.yml 

$ docker-compose up
```

## OGC Compliance testing status

GOST is being tested against the OGC SensorThings API Test Suite [https://github.com/opengeospatial/ets-sta10](https://github.com/opengeospatial/ets-sta10)

| Conformance Class                     | Reference | Implementation status |Test Status               |
|---------------------------------------|-----------|-----------------------|--------------------------| 
| Sensing Core                          | A.1       | beta                  | 6 passed, 0 failed       |
| Filtering Extension                   | A.2       | alpha                 | Testing not started      |
| Create-Update-Delete                  | A.3       | beta                  | 9 passed, 0 failed       |
| Batch Request                         | A.4       | -                     | Tests not implemented    |
| Sensing MultiDatastream Extension     | A.5       | -                     | Tests not implemented    |
| Sensing Data Array Extension          | A.6       | -                     | Tests not implemented    |
| MQTT Extension for Create and Update  | A.7       | alpha                 | Tests not implemented    |
| MQTT Extension for Receiving Updates  | A.8       | alpha                 | Tests not implemented    |

Status GOST on OGC site: [http://www.opengeospatial.org/resource/products/details/?pid=1419](http://www.opengeospatial.org/resource/products/details/?pid=1419)

## Installation and configuration

[GOST installation](docs/gost_installation.md)

[GOST configuration](docs/gost_configuration.md)

## Samples

HTTP Api: For sample requests (setting up sensors/datastreams/things and adding observations) see the tests in the [playground](test/playground_tests.md). 
For a complete collection of working requests install Postman and import the [Postman file](test/GOST.json.postman_collection) 

MQTT: For getting started with Gost and MQTT for publishing/receiving data see [GOST and MQTT - Getting started](docs/gost_mqtt_getting_started.md)

## Goals

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
