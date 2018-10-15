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
	os.Setenv("GOST_SERVER_NAME", server)
	os.Setenv("GOST_SERVER_HOST", host)
	os.Setenv("GOST_SERVER_PORT", port)
	os.Setenv("GOST_SERVER_EXTERNAL_URI", externalUri)
	os.Setenv("GOST_SERVER_MAX_ENTITIES", maxEntities)
	os.Setenv("GOST_SERVER_INDENT_JSON", indentJSON)
	os.Setenv("GOST_SERVER_HTTPS", https)
	os.Setenv("GOST_SERVER_HTTPS_KEY", httpsKey)
	os.Setenv("GOST_SERVER_HTTPS_CERT", httpsCert)
	os.Setenv("GOST_MQTT_ENABLED", mqttEnabled)
	os.Setenv("GOST_MQTT_HOST", mqttHost)
	os.Setenv("GOST_MQTT_PORT", mqttPort)
	os.Setenv("GOST_DB_HOST", dbHost)
	os.Setenv("GOST_DB_PORT", dbPort)
	os.Setenv("GOST_DB_USER", dbUser)
	os.Setenv("GOST_DB_PASSWORD", dbPassword)
	os.Setenv("GOST_DB_DATABASE", dbDB)
	os.Setenv("GOST_DB_SCHEMA", dbSchema)
	os.Setenv("GOST_DB_SSL_ENABLED", dbSSLEnabled)
	os.Setenv("GOST_DB_MAX_IDLE_CONS", dbMaxIdleCons)
	os.Setenv("GOST_DB_MAX_OPEN_CONS", dbMaxOpenCons)

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
