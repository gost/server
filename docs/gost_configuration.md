## GOST Configuration

Default configuration file: config.yaml

server: <br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;name: GOST Server (name of the webserver)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;host: localhost (host of webserver, set to 0.0.0.0 if hosting on external machine)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;port: 8080 (port of webserver)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;externalUri: http://localhost:8080/ (change to the uri where users can reach the service)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;maxEntityResponse: 50 (Max entities to return if no $top and $skip is given, not implemented yet)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;indentedJson: true (return indented JSON, not implemented yet)<br />
database:<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;host: localhost (location of PostGIS server)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;port: 5432 (port of PostGIS database)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;user: postgres (PostGIS user)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;password: postgres (PostGIS password)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;database: gost (PostGIS database to use)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;schema: v1 (schema to use)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;ssl: false (SSL enabled, not implemented yet)<br />
mqtt:<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;enabled: true (enable MQTT)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;host: iot.eclipse.org (host of the MQTT broker)<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;port: 1883 (port of the MQTT broker)<br />

The following configuration parameters can be overruled 
from the following environment variables:

db: gost_db_host, gost_db_database, gost_db_port, gost_db_user, gost_db_password. 

mqtt: gost_mqtt_host, gost_mqtt_port

server: gost_server_host, gost_server_port, gost_server_external_uri, gost_client_content

Example setting GOST environment variable on Windows:

```sh
set gost_db_host=192.168.40.10
```

Example setting GOST environment variable on Mac/Linux:

```sh
export gost_db_host=192.168.40.10
```

