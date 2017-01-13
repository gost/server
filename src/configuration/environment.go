package configuration

import (
	"os"
	"strconv"
)

// SetEnvironmentVariables changes config settings when certain environment variables are found
func SetEnvironmentVariables(conf *Config) {

	// Server settings
	gostServerName := os.Getenv("gost_server_name")
	if gostServerName != "" {
		conf.Server.Name = gostServerName
	}

	gostServerHost := os.Getenv("gost_server_host")
	if gostServerHost != "" {
		conf.Server.Host = gostServerHost
	}

	gostServerPort := os.Getenv("gost_server_port")
	if gostServerPort != "" {
		port, err := strconv.Atoi(gostServerPort)
		if err == nil {
			conf.Server.Port = int(port)
		}
	}

	gostServerExternalURI := os.Getenv("gost_server_external_uri")
	if gostServerExternalURI != "" {
		conf.Server.ExternalURI = gostServerExternalURI
	}

	gostClientContent := os.Getenv("gost_client_content")
	if gostClientContent != "" {
		conf.Server.ClientContent = gostClientContent
	}

	gostServerMaxEntities := os.Getenv("gost_server_max_entities")
	if gostServerMaxEntities != "" {
		port, err := strconv.Atoi(gostServerMaxEntities)
		if err == nil {
			conf.Server.MaxEntityResponse = int(port)
		}
	}

	gostServerIndentJSON := os.Getenv("gost_server_indent_json")
	if gostServerIndentJSON != "" {
		h, err := strconv.ParseBool(gostServerIndentJSON)
		if err == nil {
			conf.Server.IndentedJSON = h
		}
	}

	gostServerHTTPS := os.Getenv("gost_server_https")
	if gostServerHTTPS != "" {
		h, err := strconv.ParseBool(gostServerHTTPS)
		if err == nil {
			conf.Server.HTTPS = h
		}
	}

	gostServerHTTPSKey := os.Getenv("gost_server_https_key")
	if gostServerHTTPSKey != "" {
		conf.Server.HTTPSKey = gostServerHTTPSKey
	}

	gostServerHTTPSCert := os.Getenv("gost_server_https_cert")
	if gostServerHTTPSCert != "" {
		conf.Server.HTTPSCert = gostServerHTTPSCert
	}

	// MQTT settings
	gostMQTTEnabled := os.Getenv("gost_mqtt_enabled")
	if gostMQTTEnabled != "" {
		h, err := strconv.ParseBool(gostMQTTEnabled)
		if err == nil {
			conf.MQTT.Enabled = h
		}
	}

	gostMQTTHost := os.Getenv("gost_mqtt_host")
	if gostMQTTHost != "" {
		conf.MQTT.Host = gostMQTTHost
	}

	gostMQTTPort := os.Getenv("gost_mqtt_port")
	if gostMQTTPort != "" {
		port, err := strconv.Atoi(gostMQTTPort)
		if err == nil {
			conf.MQTT.Port = int(port)
		}
	}

	// Database settings
	gostDbHost := os.Getenv("gost_db_host")
	if gostDbHost != "" {
		conf.Database.Host = gostDbHost
	}

	gostDbPort := os.Getenv("gost_db_port")
	if gostDbPort != "" {
		port, err := strconv.Atoi(gostDbPort)
		if err == nil {
			conf.Database.Port = int(port)
		}
	}

	gostDbUser := os.Getenv("gost_db_user")
	if gostDbUser != "" {
		conf.Database.User = gostDbUser
	}

	gostDbPassword := os.Getenv("gost_db_password")
	if gostDbPassword != "" {
		conf.Database.Password = gostDbPassword
	}

	gostDbDatabase := os.Getenv("gost_db_database")
	if gostDbDatabase != "" {
		conf.Database.Database = gostDbDatabase
	}

	gostDbSchema := os.Getenv("gost_db_schema")
	if gostDbSchema != "" {
		conf.Database.Schema = gostDbSchema
	}

	gostDbSSL := os.Getenv("gost_mqtt_ssl_enabled")
	if gostDbSSL != "" {
		h, err := strconv.ParseBool(gostDbSSL)
		if err == nil {
			conf.Database.SSL = h
		}
	}

	gostDbMaxIdle := os.Getenv("gost_db_max_idle_cons")
	if gostDbMaxIdle != "" {
		idle, err := strconv.Atoi(gostDbMaxIdle)
		if err == nil {
			conf.Database.MaxIdleConns = idle
		}
	}

	gostDbMaxCons := os.Getenv("gost_db_max_open_cons")
	if gostDbMaxCons != "" {
		open, err := strconv.Atoi(gostDbMaxCons)
		if err == nil {
			conf.Database.MaxOpenConns = open
		}
	}
}
