
### Gost installation 


## Build from Source

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

## Build from install script

Use  [scripts/ubuntu_install.sh](scripts/ubuntu_install.sh) to install and run the latest version of GOST (including dependencies) - Tested on a clean Ubuntu 16.04 LTS installation.
