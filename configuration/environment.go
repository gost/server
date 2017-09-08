package configuration

import (
	"log"
	"os"
	"strconv"
)

// SetEnvironmentVariables changes config settings when certain environment variables are found
func SetEnvironmentVariables(conf *Config) {
	setEnvironmentServerSettings(conf)
	setEnvironmentDatabaseSettings(conf)
	setEnvironmentMQTTSettings(conf)
	setEnvironmentLoggerSettings(conf)
}

func setEnvironmentServerSettings(conf *Config) {
	gostServerName := os.Getenv("GOST_SERVER_NAME")
	if gostServerName != "" {
		conf.Server.Name = gostServerName
	}

	gostServerHost := os.Getenv("GOST_SERVER_HOST")
	if gostServerHost != "" {
		conf.Server.Host = gostServerHost
	}

	gostServerPort := os.Getenv("GOST_SERVER_PORT")
	if gostServerPort != "" {
		port, err := strconv.Atoi(gostServerPort)
		if err == nil {
			conf.Server.Port = int(port)
		}
	}

	gostServerExternalURI := os.Getenv("GOST_SERVER_EXTERNAL_URI")
	if gostServerExternalURI != "" {
		conf.Server.ExternalURI = gostServerExternalURI
		log.Println("External uri environment variable discovered: " + conf.Server.ExternalURI)

	}

	gostServerMaxEntities := os.Getenv("GOST_SERVER_MAX_ENTITIES")
	if gostServerMaxEntities != "" {
		port, err := strconv.Atoi(gostServerMaxEntities)
		if err == nil {
			conf.Server.MaxEntityResponse = int(port)
		}
	}

	gostServerIndentJSON := os.Getenv("GOST_SERVER_INDENT_JSON")
	if gostServerIndentJSON != "" {
		h, err := strconv.ParseBool(gostServerIndentJSON)
		if err == nil {
			conf.Server.IndentedJSON = h
		}
	}

	gostServerHTTPS := os.Getenv("GOST_SERVER_HTTPS")
	if gostServerHTTPS != "" {
		h, err := strconv.ParseBool(gostServerHTTPS)
		if err == nil {
			conf.Server.HTTPS = h
		}
	}

	gostServerHTTPSKey := os.Getenv("GOST_SERVER_HTTPS_KEY")
	if gostServerHTTPSKey != "" {
		conf.Server.HTTPSKey = gostServerHTTPSKey
	}

	gostServerHTTPSCert := os.Getenv("GOST_SERVER_HTTPS_CERT")
	if gostServerHTTPSCert != "" {
		conf.Server.HTTPSCert = gostServerHTTPSCert
	}
}

func setEnvironmentDatabaseSettings(conf *Config) {
	gostDbHost := os.Getenv("GOST_DB_HOST")
	if gostDbHost != "" {
		conf.Database.Host = gostDbHost
	}

	gostDbPort := os.Getenv("GOST_DB_PORT")
	if gostDbPort != "" {
		port, err := strconv.Atoi(gostDbPort)
		if err == nil {
			conf.Database.Port = int(port)
		}
	}

	gostDbUser := os.Getenv("GOST_DB_USER")
	if gostDbUser != "" {
		conf.Database.User = gostDbUser
	}

	gostDbPassword := os.Getenv("GOST_DB_PASSWORD")
	if gostDbPassword != "" {
		conf.Database.Password = gostDbPassword
	}

	gostDbDatabase := os.Getenv("GOST_DB_DATABASE")
	if gostDbDatabase != "" {
		conf.Database.Database = gostDbDatabase
	}

	gostDbSchema := os.Getenv("GOST_DB_SCHEMA")
	if gostDbSchema != "" {
		conf.Database.Schema = gostDbSchema
	}

	gostDbSSL := os.Getenv("GOST_MQTT_SSL_ENABLED")
	if gostDbSSL != "" {
		h, err := strconv.ParseBool(gostDbSSL)
		if err == nil {
			conf.Database.SSL = h
		}
	}

	gostDbMaxIdle := os.Getenv("GOST_DB_MAX_IDLE_CONS")
	if gostDbMaxIdle != "" {
		idle, err := strconv.Atoi(gostDbMaxIdle)
		if err == nil {
			conf.Database.MaxIdleConns = idle
		}
	}

	gostDbMaxCons := os.Getenv("GOST_DB_MAX_OPEN_CONS")
	if gostDbMaxCons != "" {
		open, err := strconv.Atoi(gostDbMaxCons)
		if err == nil {
			conf.Database.MaxOpenConns = open
		}
	}
}

func setEnvironmentMQTTSettings(conf *Config) {
	gostMQTTEnabled := os.Getenv("GOST_MQTT_ENABLED")
	if gostMQTTEnabled != "" {
		h, err := strconv.ParseBool(gostMQTTEnabled)
		if err == nil {
			conf.MQTT.Enabled = h
		}
	}

	gostMQTTHost := os.Getenv("GOST_MQTT_HOST")
	if gostMQTTHost != "" {
		conf.MQTT.Host = gostMQTTHost
	}

	gostMQTTPort := os.Getenv("GOST_MQTT_PORT")
	if gostMQTTPort != "" {
		port, err := strconv.Atoi(gostMQTTPort)
		if err == nil {
			conf.MQTT.Port = int(port)
		}
	}
}

func setEnvironmentLoggerSettings(conf *Config) {
	gostLoggerFileName := os.Getenv("GOST_LOG_FILENAME")
	if gostLoggerFileName != "" {
		conf.Logger.FileName = gostLoggerFileName
	}

	gostLoggerVerbose := os.Getenv("GOST_LOG_VERBOSE_FLAG")
	if gostLoggerVerbose != "" {
		if verboseFlag, err := strconv.ParseBool(gostLoggerVerbose); err == nil {
			conf.Logger.Verbose = verboseFlag
		}
	}
}
