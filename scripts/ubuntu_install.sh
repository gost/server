#-------------------------
# update system and tools
#-------------------------
sudo apt-get -y update
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


#Open port
sudo iptables -A INPUT -m state --state NEW -m tcp -p tcp --dport 8080 -j ACCEPT -m comment --comment "GOST Server port"

