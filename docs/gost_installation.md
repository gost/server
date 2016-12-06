
### Gost installation 

In this document 4 ways to install GOST server are described: Installation from pre-build binaries (for Linux/Windows/Mac), 
Installation from source code (requires some Golang knowledge), installing with Docker amd complete GOST installation from script. 

## Prerequisites

1) Install PostgreSQL (http://www.postgresql.org/) with PostGIS and GOST database

See https://github.com/Geodan/gost-db/blob/master/README.md for various options

Testcase after installation: 

In browser open http://localhost:8080/v1.0 to test if the server is running

## Install from Binaries

todo

## Install from Source

1) Install GoLang

See https://golang.org/

2) Clone code
```sh
git clone https://github.com/Geodan/gost.git
```
3) Get dependencies

```sh
go get github.com/gorilla/mux
go get gopkg.in/yaml.v2
go get github.com/lib/pq
go get github.com/eclipse/paho.mqtt.golang
```

4) Edit config.yaml or set environment settings to change connection to database<br />

5) Start

```sh
go run main.go
```

### Install with Docker

```
$ wget https://raw.githubusercontent.com/Geodan/gost/master/src/docker-compose.yml 

$ docker-compose up
```

## Build from install script

For a complete GOST server installation from script there are various parts to install

1) Install database

Run script at https://github.com/Geodan/gost-db/blob/master/gost-db-install.sh for installing Postgres + Postgis + GOST database

2) Install Goland

// todo

2) Install GOST server

// todo

3) Install Mosquitto

// todo

4) Configure system: opening ports

// todo

5) Configure system: set services resatrt options

// todo

See  [../scripts/ubuntu_install.sh](../scripts/ubuntu_install.sh) to for some examples how to get this working on your system.
