## GOST Docker support

For getting GOST to work in Docker there are two images available on Docker Hub:

For service and dashboard: [https://hub.docker.com/r/geodan/gost/]

For database: [https://hub.docker.com/r/geodan/gost-db/]

The docker images can run separately, or running in a combined way using the Dockercompose file. 

If you want to use the latest (unstable) build of GOST, you can use the automated builds found here: [https://hub.docker.com/r/geodan/gost-nightly/] and here: [https://hub.docker.com/r/geodan/gost-db-nightly/] there is no guarantee these images will run correctly. 

# Running GOST with Docker-compose

```
$ wget https://raw.githubusercontent.com/Geodan/gost/master/src/docker-compose.yml 

$ docker-compose up
```

# Running GOST database

```
. $ docker run -p 5432:5432 --name -e POSTGRES_DB=gost geodan/gost-db
```

Connect in pgadmin with localhost:5432 postgres/postgres

GOST schema is in schema postgres.v1

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
# Building GOST service and dashboard image

```
. $ git clone https://github.com/Geodan/gost.git

. $ cd src/github.com/geodan/gost/src

. $ docker build -t geodan/gost:latest .

. $ docker push geodan/gost

```

# Building GOST-db image

```
. $ git clone https://github.com/Geodan/gost.git

. $ cd src/github.com/geodan/gost/src/docker/postgis

. $ docker build -t geodan/gost-db .

. $ docker push geodan/gost-db
```
