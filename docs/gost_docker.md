## Docker support

Work has started on Docker support for GOST. The GOST Docker image is available at
[https://hub.docker.com/r/geodan/gost/] and is automatic rebuild after a Github commit.

# Running GOST database

. $ docker run -p 5432:5432 -d geodan/gost-db

Connect in pgadmin with localhost:5432 postgres/postgres

GOST schema is in schema postgres.v1

# Running GOST service and dashboard

$ docker run -p 8080:8080 -t geodan/gost

GOST is available at http://localhost:8080 

TODO: describe connection to database

# Running GOST with Docker-compose

. Clone the repository

. $ cd src/github.com/geodan/gost/src

. $ docker-compose up

# Building GOST service and dashboard image

. Clone the repository

. $ cd src/github.com/geodan/gost/src

. $ docker build -t geodan/gost:latest .

# Building GOST-db image

. Clone the repository

. $ cd src/github.com/geodan/gost/src/docker/postgis

. $ docker build -t geodan/gost-db .
