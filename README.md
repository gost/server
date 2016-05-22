<img src="src/client/assets/img/icon.png" width="353"><br />
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Geodan/gost/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/Geodan/gost?status.svg)](https://godoc.org/github.com/Geodan/gost)
[![Build Status](http://beta.drone.io/api/badges/drone/drone/status.svg)](https://drone.io/github.com/Geodan/gost/latest)
[![Go Report Card](https://goreportcard.com/badge/geodan/gost)](https://goreportcard.com/report/geodan/gost)
[![Coverage Status](https://coveralls.io/repos/github/Geodan/gost/badge.svg?branch=master)](https://coveralls.io/github/Geodan/gost?branch=master)
[![Join the chat at https://gitter.im/Geodan/gost](https://badges.gitter.im/Geodan/gost.svg)](https://gitter.im/Geodan/gost?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)<br />

GOST (Go-SensorThings) is a sensor server written in Go. It implements the [OGC SensorThings API] (http://ogc-iot.github.io/ogc-iot-api/api.html) standard.

## Disclaimer

GOST is alpha software and is not considered appropriate for customer use. Feel free to help development :-)

## Binaries

Release 0.1: 2016-05-11 [https://github.com/Geodan/gost/releases/tag/0.1](https://github.com/Geodan/gost/releases/tag/0.1)
Binaries are build for Windows, Ubuntu and OSX.
A cross-compilation file can be found here [scripts/xcompile.bat](https://github.com/Geodan/gost/blob/master/scripts/xcompile.bat)

## Getting started for developers

1) Install GoLang (https://golang.org/)<br />
2) Install Postgresql (http://www.postgresql.org/) and PostGIS <br />
Create a database and run "CREATE EXTENSION postgis;" on it<br />
3) Clone code
```sh
git clone https://github.com/Geodan/gost.git
```
4) Get dependencies
```sh
go get github.com/gorilla/mux
go get gopkg.in/yaml.v2
go get github.com/lib/pq
go get github.com/eclipse/paho.mqtt.golang
```
5) Edit config.yaml or set environment settings to change connection to database<br />
6) Create GOST schema
```sh
go run main.go -install ../scripts/createdb.sql
```
7) Start
```sh
go run main.go
```

8) In browser open http://localhost:8080/v1.0 to test if the server is running

<br /><br />
Or try [scripts/ubuntu_install.sh](scripts/ubuntu_install.sh) to install and run the latest version of GOST (including dependencies) - Tested on a clean Ubuntu 16.04 LTS installation.

## Sample requests

For sample requests (setting up sensors/datastreams/things and adding observations) see the tests in the [playground](test/playground_tests.md). 
For a complete collection of working requests install Postman and import the [Postman file](test/GOST.json.postman_collection) 

## GOST Startup flags

-config "path to config file": specify the config file (default config.yaml)<br />
-install "path to database creation file": creates the database schema

## Configuration

Default file: config.yaml

server: <br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;name: GOST Server (name of the webserver)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;host: localhost (host of webserver, set to 0.0.0.0 if hosting on external machine)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;port: 8080 (port of webserver)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;externalUri: http://localhost:8080/ (change to the uri where users can reach the service)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;clientContent: ./client/ (Location of the client folder (dashboard))<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;maxEntityResponse: 50 (Max entities to return if no $top and $skip is given, not implemented yet)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;indentedJson: true (return indented JSON, not implemented yet)<br />
database:<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;host: localhost (location of PostGIS server)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;port: 5432 (port of PostGIS database)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;user: postgres (PostGIS user)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;password: postgres (PostGIS password)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;database: gost (PostGIS database to use)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;schema: v1 (schema to use)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;ssl: false (SSL enabled, not implemented yet)<br />
mqtt:<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;enabled: true (enable MQTT)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;host: iot.eclipse.org (host of the MQTT broker)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;port: 1883 (port of the MQTT broker)<br />

The following configuration parameters can be overruled 
from the following environment variables:

db: gost_db_host, gost_db_port, gost_db_user, gost_db_password. 

mqtt: gost_mqtt_host, gost_mqtt_port

server: gost_server_host, gost_server_port, gost_server_external_uri, gost_client_content

Example setting Gost environment variable on Windows:

```sh
set gost_db_host=192.168.40.10
```

Example setting Gost environment variable on Mac/Linux:

```sh
export gost_db_host=192.168.40.10
```

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
