#-------------------------
# update system and tools
# Note: Mosquitto installed trough apt-get does not support websockets
# to get websockets working check: https://github.com/Geodan/gost/wiki/Mosquitto-with-websockets
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

#-------------------------
# Port configuration
#-------------------------
sudo iptables -A INPUT -m state --state NEW -m tcp -p tcp --dport 80 -j ACCEPT -m comment --comment "GOST Server port"
sudo iptables -A INPUT -m state --state NEW -m tcp -p tcp --dport 1883 -j ACCEPT -m comment --comment "Mosquitto MQTT port"

#Add port to firewall
sudo ufw allow 1883
sudo ufw allow 80

#-------------------------
# Get latest version of GOST from github
#-------------------------
cd ~/dev/go/src/github.com/geodan
git clone https://github.com/Geodan/gost.git
cd gost/src
go get .

#-------------------------
# Build GOST to bin folder
#-------------------------
sudo mkdir /usr/local/bin/gost
go build -o /usr/local/bin/gost github.com/geodan/gost/src

#-------------------------
# Copy needed files to bin folder
#-------------------------
sudo cp ~/dev/go/src/github.com/geodan/gost/scripts/createdb.sql /usr/local/bin/gost
sudo cp -avr ~/dev/go/src/github.com/geodan/gost/src/client /usr/local/bin/gost

#Create schema in Postgresql
/usr/local/bin/gost/gost -install /usr/local/bin/gost/config.yaml

#-------------------------
# Create /etc/systemd/system/gost.service to run GOST as a service
#-------------------------
echo "
[Unit]
Description=GOST Server
After=syslog.target network.target postgresql.service

[Service]
Environment=export gost_server_host=0.0.0.0
Environment=export gost_server_port=80
Environment=export gost_server_external_uri=http://mysite.com
Environment=export gost_server_client_content=/usr/local/bin/gost/client/
ExecStart=/usr/local/bin/gost/gost -config /usr/local/bin/gost/config.yaml

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/gost.service

#-------------------------
# Enable GOST service on boot, start GOST
#-------------------------
sudo systemctl daemon-reload
sudo systemctl enable gost
sudo systemctl start gost