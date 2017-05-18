package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"os"
)

func TestEnvironmentVariabels(t *testing.T) {
	// arrange
	conf := Config{}

	// act
	os.Setenv("gost_server_name", "server")
	os.Setenv("gost_server_host", "host")
	os.Setenv("gost_server_port", "8080")
	os.Setenv("gost_server_external_uri", "uri")
	os.Setenv("gost_client_content", "client_content")
	os.Setenv("gost_server_max_entities", "1")
	os.Setenv("gost_server_indent_json", "true")
	os.Setenv("gost_server_https", "false")
	os.Setenv("gost_server_https_key", "key")
	os.Setenv("gost_server_https_cert", "cert")
	os.Setenv("gost_mqtt_enabled", "true")
	os.Setenv("gost_mqtt_host", "mqtt_host")
	os.Setenv("gost_mqtt_port", "9001")
	os.Setenv("gost_db_host", "db_host")
	os.Setenv("gost_db_port", "5432")
	os.Setenv("gost_db_user", "db_user")
	os.Setenv("gost_db_password", "db_password")
	os.Setenv("gost_db_database", "db_database")
	os.Setenv("gost_db_schema", "db_schema")
	os.Setenv("gost_mqtt_ssl_enabled", "true")
	os.Setenv("gost_db_max_idle_cons", "1")
	os.Setenv("gost_db_max_open_cons", "1")

	SetEnvironmentVariables(&conf)

	// assert
	assert.NotNil(t, conf, "Configuration should not be nil")
}
