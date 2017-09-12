<img src="https://raw.githubusercontent.com/gost/dashboard/master/content/assets/img/icon.png" width="353"><br />
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/gost/server/blob/master/LICENSE)
[![API DOC](https://img.shields.io/badge/api%20doc-reference-blue.svg)](http://docs.gost1.apiary.io/#)
[![GoDoc](https://godoc.org/github.com/gost/server?status.svg)](https://godoc.org/github.com/gost/server)
[![Build Status](https://travis-ci.org/gost/server.svg?branch=master)](https://travis-ci.org/gost/server)
[![Go Report Card](https://goreportcard.com/badge/gost/server)](https://goreportcard.com/report/gost/server)
[![Coverage Status](https://coveralls.io/repos/github/gost/server/badge.svg?branch=master)](https://coveralls.io/github/gost/server?branch=master)
[![Join the chat at https://gitter.im/gost/Lobby](https://badges.gitter.im/gost/gost.svg)](https://gitter.im/gost/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Join Slack](https://slackinvitergost.herokuapp.com/badge.svg)](https://slackinvitergost.herokuapp.com/)<br />

GOST (Go-SensorThings) is an IoT Platform written in Golang (Go). It implements the Sensing profile (part 1) of the [OGC SensorThings API] (http://ogc-iot.github.io/ogc-iot-api/api.html) standard including the MQTT extension.

<a href="http://www.opengeospatial.org/resource/products/details/?pid=1419"><img src ="https://raw.githubusercontent.com/gost/docs/master/images/Certified_OGC_Compliant_Logo_Web.png"/><br/></a>

Implementation of the Tasking profile (part 2) and Rules Engine profile (part 3) of the OGC SensorThings API is planned as a future work activity.

The GOST website and blog can be found at [www.gostserver.xyz](https://www.gostserver.xyz)

## Disclaimer

GOST is alpha software and is not (yet) considered appropriate for customer use. Feel free to help development :-)

## Binaries

Current release [V0.5](https://github.com/gost/server/releases/tag/0.5): 2017-07-17

Binaries are build for Windows, Ubuntu and OSX.


## Roadmap

| Date       	|             Version 	| Features                                                        	|
|------------	|---------------------	|-----------------------------------------------------------------	|
| 2017-09-17 	| 0.6                 	| TBD                                                             	|


## Running GOST server with docker run

- Start geodan/gost-db which creates an user postgres with password postgres and initialises a database named gost
- Start geodan/gost and set info to connect to gost-db, gost will be available at http://localhost:8080/v1.0
- Start geodan/gost-dashboard and link gost to use GOST trough nginx, gost + dashboard available at http://localhost:8081

```
$ docker run -d -p 5432:5432 -e POSTGRES_DB=gost -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres --name gost-db geodan/gost-db
$ docker run -d -p 8080:8080 --link gost-db:gost-db -e GOST_DB_HOST=gost-db -e GOST_DB_USER=postgres -e GOST_DB_PASSWORD=postgres --name gost geodan/gost
$ docker run -d -p 8081:8080 --link gost:gost --name gost-dashboard geodan/gost-dashboard		
```

For making connection to external database use environmental variables GOST_DB_HOST, GOST_DB_PORT, GOST_DB_DATABASE, GOST_DB_USER, GOST_DB_PASSWORD
```
$ docker run -d -p 8080:8080 -t -e GOST_DB_HOST=192.168.40.10 -e GOST_DB_DTABASE=gost --name gost geodan/gost
```

If you want to load your own config file, mount a location with your config.yaml file to /gostserver/config. If the filename does not equals config.yaml you have to run gost with the -config parameter:
```
$ docker run -v myconfiglocation:/gostserver/config geodan/gost -config /gostserver/config/myconfig.yaml
```

## Building GOST server 

```
$ git clone https://github.com/Geodan/gost.git

$ cd src/github.com/geodan/gost/src

$ docker build -t geodan/gost:latest .

$ docker push geodan/gost

```
## Building GOST server for Raspberry Pi

note: building the Raspberry Pi image must be done on a Raspberry Pi :-(, otherwise errors will occur.

```
pi@raspberrypi:~/dev/go/src/github.com/Geodan/gost $ sudo docker build -f Dockerfile-rpi -t geodan/rpi-gost .

pi@raspberrypi:~/dev/go/src/github.com/Geodan/gost $ sudo docker push geodan/rpi-gost
```

## OGC Compliance testing status

GOST is being tested against the OGC SensorThings API Test Suite 1.0 [https://github.com/opengeospatial/ets-sta10](https://github.com/opengeospatial/ets-sta10)

| Conformance Class                     | Reference | Implementation status |Test Status               |
|---------------------------------------|-----------|-----------------------|--------------------------| 
| Sensing Core                          | A.1       | beta                  | 6/6                      |
| Filtering Extension                   | A.2       | beta                  | 8/8                      |
| Create-Update-Delete                  | A.3       | beta                  | 9/9                      |
| Batch Request                         | A.4       | -                     | Tests not implemented    |
| Sensing MultiDatastream Extension     | A.5       | -                     | Tests not implemented    |
| Sensing Data Array Extension          | A.6       | -                     | Tests not implemented    |
| MQTT Extension for Create and Update  | A.7       | alpha                 | Tests not implemented    |
| MQTT Extension for Receiving Updates  | A.8       | alpha                 | Tests not implemented    |

Status GOST on OGC site: [http://www.opengeospatial.org/resource/products/details/?pid=1419](http://www.opengeospatial.org/resource/products/details/?pid=1419)

## Installation and configuration

[GOST installation](https://github.com/gost/docs/blob/master/gost_installation.md)

[GOST configuration](https://github.com/gost/docs/blob/master/gost_configuration.md)

## Security

[GOST security](https://github.com/gost/docs/blob/master/gost_security.md)

## Samples
[Apiary API Docs](http://docs.gost1.apiary.io/)  

HTTP API: For sample requests (setting up sensors/datastreams/things and adding observations) see the tests in the [playground](https://github.com/gost/docs/blob/master/playground_tests.md). 
For a complete collection of working requests install Postman and import the [Postman file](https://github.com/gost/postman/blob/master/GOST.postman_collection.json) 

MQTT: For getting started with Gost and MQTT for publishing/receiving data see [GOST and MQTT - Getting started](https://github.com/gost/docs/blob/master/gost_mqtt_getting_started.md)

## Goals

- Complete implementation of the OGC SensorThings spec
- Test coverage!
- Frontend
- Benchmarks
- Authentication
- Different storage providers such as MongoDB (Now using PostgreSQL)

## Benchmarks

[Publish observation (MQTT)](https://github.com/Geodan/gost/wiki/Benchmark---publish-observation-(MQTT))
