#-------------------------
# update system and tools
#-------------------------
sudo apt -y update
sudo apt-get -y install git
sudo apt-get -y install mosquitto

#-------------------------
# create dirs
#-------------------------
cd ~
mkdir -p dev/go/src/github.com/geodan

#-------------------------
# install golang
#-------------------------
sudo apt-get -y install golang
export GOPATH=$HOME/dev/go
export PATH=$PATH:$GOPATH/bin

#-------------------------
# install postgresql + postgis
#-------------------------
sudo apt-get -y install postgresql postgresql-contrib postgis
sudo su postgres -c psql << EOF
ALTER USER postgres WITH PASSWORD 'postgres';
CREATE DATABASE gost OWNER postgres;
\connect gost
CREATE EXTENSION postgis;
\q
EOF

#-------------------------
# install gost
#-------------------------
cd ~/dev/go/src/github.com/geodan
git clone https://github.com/Geodan/gost.git
go get github.com/gorilla/mux
go get gopkg.in/yaml.v2
go get github.com/lib/pq
go get github.com/eclipse/paho.mqtt.golang

#Build GOST to bin folder
go build -o gost/bin/gost github.com/geodan/gost/src

#Copy the createdb script to bin folder
sudo cp ~/dev/go/src/github.com/geodan/gost/scripts/createdb.sql ~/dev/go/src/github.com/geodan/gost/bin

#Copy the client folder to bin
sudo cp -avr ~/dev/go/src/github.com/geodan/gost/src/client ~/dev/go/src/github.com/geodan/gost/bin

#Write a config.yaml
echo "server:
    name: GOST Server
    host: 0.0.0.0
    port: 8080
    externalUri: http://localhost:8080/
database:
    host: localhost
    port: 5432
    user: postgres
    password: postgres
    database: gost
    schema: v1
    ssl: false
mqtt:
    enabled: true
    host: 127.0.0.1
    port: 1883" > ~/dev/go/src/github.com/geodan/gost/bin/config.yaml

#Setup database schema for GOST
~/dev/go/src/github.com/geodan/gost/bin/gost -install ~/dev/go/src/github.com/geodan/gost/bin/createdb.sql -config ~/dev/go/src/github.com/geodan/gost/bin/config.yaml

#Run GOST
~/dev/go/src/github.com/geodan/gost/bin/gost -config ~/dev/go/src/github.com/geodan/gost/bin/config.yaml