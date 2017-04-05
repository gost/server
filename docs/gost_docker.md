## GOST Docker support

For getting GOST to work in Docker there are two images available on Docker Hub:

For service and dashboard: [https://hub.docker.com/r/geodan/gost/]

For database: [https://hub.docker.com/r/geodan/gost-db/]

For more information about the Docker gost-db image, see [https://github.com/Geodan/gost-db]

The docker images can run separately, or running in a combined way using the Dockercompose file.

Tags: Use the tag latest for the latest development version, otherwise use a tag like '0.4' for more stable versions.

# Running GOST with Docker-compose

```
$ wget https://raw.githubusercontent.com/Geodan/gost/master/src/docker-compose.yml 

$ docker-compose up
```

# Running GOST service and dashboard
```
$ docker run -p 8080:8080 --link gost-db:gost-db -e gost_db_host=gost-db geodan/gost
```
GOST is available at http://localhost:8080 

For making connection to external database use environmental variables gost_db_host, gost_db_port, gost_db_user, gost_db_password:

For example: 
```
docker run -p 8080:8080 -t -e gost_db_host=192.168.40.10 -e gost_db_database=gost geodan/gost

```

on raspberrypi:
```
docker run -p 8080:8080 -t -e gost_db_host=raspberrypi -e gost_db_database=gost -e gost_mqtt_host=raspberrypi geodan/rpi-gost
```

# Building GOST service and dashboard image

```
. $ git clone https://github.com/Geodan/gost.git

. $ cd src/github.com/geodan/gost/src

. $ docker build -t geodan/gost:latest .

. $ docker push geodan/gost

```
# Building GOST service and dashboard image for Raspberrypi

```
docker build -f Dockerfile-rpi -t geodan/rpi-gost .

docker push geodan/rpi-gost
```

