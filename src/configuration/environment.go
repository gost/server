package configuration

import (
	"os"
	"strconv"
)

// SetEnvironmentVariables changes config settings when certain environment variables are found
func SetEnvironmentVariables(conf *Config) {
	gostMqttHost := os.Getenv("gost_mqtt_host")
	if gostMqttHost != "" {
		conf.MQTT.Host = gostMqttHost
	}

	gostMqttPort := os.Getenv("gost_mqtt_port")
	if gostMqttPort != "" {
		port, err := strconv.Atoi(gostMqttPort)
		if err == nil {
			conf.MQTT.Port = int(port)
		}
	}

	gostServerHost := os.Getenv("gost_server_host")
	if gostServerHost != "" {
		conf.Server.Host = gostServerHost
	}

	gostServerExternalURI := os.Getenv("gost_server_external_uri")
	if gostServerExternalURI != "" {
		conf.Server.ExternalURI = gostServerExternalURI
	}

	gostClientContent := os.Getenv("gost_client_content")
	if gostClientContent != "" {
		conf.Server.ClientContent = gostClientContent
	}

	gostServerPort := os.Getenv("gost_server_port")
	if gostServerPort != "" {
		port, err := strconv.Atoi(gostServerPort)
		if err == nil {
			conf.Server.Port = int(port)
		}
	}
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
}
