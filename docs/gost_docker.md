## Docker support

Work has started on Docker support for GOST. The GOST Docker image is available at
[https://hub.docker.com/r/geodan/gost/] and is automatic rebuild after a Github commit.

# Running GOST with Docker-compose

. Clone the repository

. $ cd src/github.com/geodan/gost/src

. $ docker-compose up

# Building GOST Docker image

. Clone the repository

. $ cd src/github.com/geodan/gost/src

. $ docker build -t geodan/gost:latest .

# Running GOST in Docker

$ docker run -p 8080:8080 -t geodan/gost:latest

Site should run on port 8080

Todo: create docker-compose file




