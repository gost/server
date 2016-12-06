
### Gost installation 

## Prerequisites

1) Install GoLang

See (https://golang.org/)

2) Install PostgreSQL (http://www.postgresql.org/) with PostGIS and GOST database

See (https://github.com/Geodan/gost-db/blob/master/README.md)

## Install from binaries

todo

## Install from Source

1) Clone code
```sh
git clone https://github.com/Geodan/gost.git
```
2) Get dependencies

```sh
go get github.com/gorilla/mux
go get gopkg.in/yaml.v2
go get github.com/lib/pq
go get github.com/eclipse/paho.mqtt.golang
```

3) Edit config.yaml or set environment settings to change connection to database<br />

4) Start

```sh
go run main.go
```

5) In browser open http://localhost:8080/v1.0 to test if the server is running

## Build from install script

Use  [../scripts/ubuntu_install.sh](../scripts/ubuntu_install.sh) to install and run the latest version of GOST (including dependencies) - Tested on a clean Ubuntu 16.04 LTS installation.
