<img src="src/client/assets/img/icon.png" width="353"><br />
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Geodan/gost/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/Geodan/gost?status.svg)](https://godoc.org/github.com/Geodan/gost)
[![Build Status](http://beta.drone.io/api/badges/drone/drone/status.svg)](https://drone.io/github.com/Geodan/gost/latest)
[![Go Report Card](https://goreportcard.com/badge/geodan/gost)](https://goreportcard.com/report/geodan/gost)
[![Coverage Status](https://coveralls.io/repos/github/Geodan/gost/badge.svg?branch=master)](https://coveralls.io/github/Geodan/gost?branch=master)
[![Join the chat at https://gitter.im/Geodan/gost](https://badges.gitter.im/Geodan/gost.svg)](https://gitter.im/Geodan/gost?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)<br />

GOST (Go-SensorThings) is a sensor server written in Go. It implements the [OGC SensorThings API] (http://ogc-iot.github.io/ogc-iot-api/api.html) standard including the MQTT extension.

## Disclaimer

GOST is alpha software and is not (yet) considered appropriate for customer use. Feel free to help development :-)

## Binaries

Release 0.1: 2016-05-11 [https://github.com/Geodan/gost/releases/tag/0.1] (https://github.com/Geodan/gost/releases/tag/0.1)

Binaries are build for Windows, Ubuntu and OSX.

## Installation and configuration

[Gost installation](docs/gost_installation.md)

[Gost configuration](docs/gost_configuration.md)

## Sample requests

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
