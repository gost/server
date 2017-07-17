package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
)

func TestEnvironmentVariabelsEmpty(t *testing.T) {
	// arrange
	conf := Config{}

	// act
	SetEnvironmentVariables(&conf)

	// assert
	assert.NotNil(t, conf, "Configuration should not be nil")
}

func TestEnvironmentVariabels(t *testing.T) {
	// arrange
	conf := Config{}
	server := "server"
	host := "host"
	port := "8080"
	portParsed, _ := strconv.Atoi(port)
	externalUri := "uri"
	maxEntities := "1"
	maxEntitiesParsed, _ := strconv.Atoi(maxEntities)
	indentJSON := "true"
	indentJSONParsed, _ := strconv.ParseBool(indentJSON)
	https := "true"
	httpsParsed, _ := strconv.ParseBool(https)
	httpsKey := "key"
	httpsCert := "cert"
	mqttEnabled := "true"
	mqttEnabledParsed, _ := strconv.ParseBool(mqttEnabled)
	mqttHost := "mqtt_host"
	mqttPort := "9001"
	mqttPortParsed, _ := strconv.Atoi(mqttPort)
	dbSSLEnabled := "true"
	dbSSLEnabledParsed, _ := strconv.ParseBool(dbSSLEnabled)
	dbHost := "db_host"
	dbPort := "5432"
	dbPortParsed, _ := strconv.Atoi(dbPort)
	dbUser := "user"
	dbPassword := "secret"
	dbDB := "gost"
	dbSchema := "v1"
	dbMaxIdleCons := "1"
	dbMaxIdleConsParsed, _ := strconv.Atoi(dbMaxIdleCons)
	dbMaxOpenCons := "1"
	dbMaxOpenConsParsed, _ := strconv.Atoi(dbMaxOpenCons)
	// act
	os.Setenv("gost_server_name", server)
	os.Setenv("gost_server_host", host)
	os.Setenv("gost_server_port", port)
	os.Setenv("gost_server_external_uri", externalUri)
	os.Setenv("gost_server_max_entities", maxEntities)
	os.Setenv("gost_server_indent_json", indentJSON)
	os.Setenv("gost_server_https", https)
	os.Setenv("gost_server_https_key", httpsKey)
	os.Setenv("gost_server_https_cert", httpsCert)
	os.Setenv("gost_mqtt_enabled", mqttEnabled)
	os.Setenv("gost_mqtt_host", mqttHost)
	os.Setenv("gost_mqtt_port", mqttPort)
	os.Setenv("gost_db_host", dbHost)
	os.Setenv("gost_db_port", dbPort)
	os.Setenv("gost_db_user", dbUser)
	os.Setenv("gost_db_password", dbPassword)
	os.Setenv("gost_db_database", dbDB)
	os.Setenv("gost_db_schema", dbSchema)
	os.Setenv("gost_mqtt_ssl_enabled", dbSSLEnabled)
	os.Setenv("gost_db_max_idle_cons", dbMaxIdleCons)
	os.Setenv("gost_db_max_open_cons", dbMaxOpenCons)

	SetEnvironmentVariables(&conf)

	// assert
	assert.NotNil(t, conf, "Configuration should not be nil")
	assert.Equal(t, server, conf.Server.Name)
	assert.Equal(t, externalUri, conf.Server.ExternalURI)
	assert.Equal(t, host, conf.Server.Host)
	assert.Equal(t, httpsParsed, conf.Server.HTTPS)
	assert.Equal(t, httpsCert, conf.Server.HTTPSCert)
	assert.Equal(t, httpsKey, conf.Server.HTTPSKey)
	assert.Equal(t, indentJSONParsed, conf.Server.IndentedJSON)
	assert.Equal(t, maxEntitiesParsed, conf.Server.MaxEntityResponse)
	assert.Equal(t, portParsed, conf.Server.Port)
	assert.Equal(t, mqttHost, conf.MQTT.Host)
	assert.Equal(t, mqttEnabledParsed, conf.MQTT.Enabled)
	assert.Equal(t, mqttPortParsed, conf.MQTT.Port)
	assert.Equal(t, dbDB, conf.Database.Database)
	assert.Equal(t, dbPortParsed, conf.Database.Port)
	assert.Equal(t, dbHost, conf.Database.Host)
	assert.Equal(t, dbMaxIdleConsParsed, conf.Database.MaxIdleConns)
	assert.Equal(t, dbMaxOpenConsParsed, conf.Database.MaxOpenConns)
	assert.Equal(t, dbPassword, conf.Database.Password)
	assert.Equal(t, dbSchema, conf.Database.Schema)
	assert.Equal(t, dbSSLEnabledParsed, conf.Database.SSL)
	assert.Equal(t, dbUser, conf.Database.User)
}
