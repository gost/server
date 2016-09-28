## Docker support

Work has started on Docker support for GOST

# Building GOST Docker image

. Clone the repository
. $ cd src/github.com/geodan/gost/src
. $ docker build -t geodan/gost:latest .

# Running GOST in Docker

$ docker run -p 5001:8080 -t geodan/gost:latest
Site should run on port 8080

Todo: create docker-compose file




