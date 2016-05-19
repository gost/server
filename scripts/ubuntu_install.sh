#-------------------------
# update system and tools
#-------------------------
sudo apt-get update
sudo apt-get -y install git
sudo apt-get -y install mosquitto

#-------------------------
# create dirs
#-------------------------
cd ~
rm -rf dev
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

#Open port
sudo iptables -A INPUT -m state --state NEW -m tcp -p tcp --dport 8080 -j ACCEPT -m comment --comment "GOST Server port"
sudo iptables -A INPUT -m state --state NEW -m tcp -p tcp --dport 1883 -j ACCEPT -m comment --comment "Mosquitto MQTT port"
sudo iptables -A INPUT -m state --state NEW -m tcp -p tcp --dport 9001 -j ACCEPT -m comment --comment "Mosquitto MQTT Websocket port"



#-------------------------
# install gost
#-------------------------
cd ~/dev/go/src/github.com/geodan
git clone https://github.com/Geodan/gost.git
cd gost/src
go get .

export gost_server_host=0.0.0.0
export gost_server_external_uri=http://37.97.183.133:8080/

go run main.go -install ../scripts/createdb.sql
go run main.go
