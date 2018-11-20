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

	gostDbSSL := os.Getenv("GOST_DB_SSL_ENABLED")
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

	gostMQTTVerbose := os.Getenv("GOST_MQTT_VERBOSE")
	if gostMQTTVerbose != "" {
		h, err := strconv.ParseBool(gostMQTTVerbose)
		if err == nil {
			conf.MQTT.Verbose = h
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
			conf.MQTT.Port = port
		}
	}

	gostMQTTPrefix := os.Getenv("GOST_MQTT_PREFIX")
	if gostMQTTPrefix != "" {
		conf.MQTT.Prefix = gostMQTTPrefix
	}

	gostMQTTClientID := os.Getenv("GOST_MQTT_CLIENTID")
	if gostMQTTClientID != "" {
		conf.MQTT.ClientID = gostMQTTClientID
	}

	gostMQTTSubscriptionQOS := os.Getenv("GOST_MQTT_SUBSCRIPTIONQOS")
	if gostMQTTSubscriptionQOS != "" {
		qos,err := strconv.ParseInt(gostMQTTSubscriptionQOS,0,8)
		if err == nil {
			conf.MQTT.SubscriptionQos = byte(qos)
		}

	}

	gostMQTTPersistent := os.Getenv("GOST_MQTT_PERSISTENT")
	if gostMQTTPersistent != "" {
		persistent, err := strconv.ParseBool(gostMQTTPersistent)
		if err == nil {
			conf.MQTT.Persistent = persistent
		}
	}

	gostMQTTOrder := os.Getenv("GOST_MQTT_ORDER_MATTERS")
	if gostMQTTOrder != "" {
		order, err := strconv.ParseBool(gostMQTTOrder)
		if err == nil {
			conf.MQTT.Order = order
		}
	}

	gostMQTTssl := os.Getenv("GOST_MQTT_SSL")
	if gostMQTTssl != "" {
		ssl, err := strconv.ParseBool(gostMQTTssl)
		if err == nil {
			conf.MQTT.SSL = ssl
		}
	}

	gostMQTTUsername := os.Getenv("GOST_MQTT_USERNAME")
	if gostMQTTUsername != "" {
		conf.MQTT.Username = gostMQTTUsername
	}

	gostMQTTPassword := os.Getenv("GOST_MQTT_PASSWORD")
	if gostMQTTPassword != "" {
		conf.MQTT.Password = gostMQTTPassword
	}

	gostMQTTCaCertFile := os.Getenv("GOST_MQTT_CA_CERT_FILE")
	if gostMQTTCaCertFile != "" {
		conf.MQTT.CaCertFile = gostMQTTCaCertFile
	}

	gostMQTTClientCertFile := os.Getenv("GOST_MQTT_CLIENT_CERT_FILE")
	if gostMQTTClientCertFile != ""{
		conf.MQTT.ClientCertFile = gostMQTTClientCertFile
	}


	gostMQTTPrivateKeyFile := os.Getenv("GOST_MQTT_PRIVATE_KEY_FILE")
	if gostMQTTPrivateKeyFile != ""{
		conf.MQTT.PrivateKeyFile = gostMQTTPrivateKeyFile
	}

	keepaliveSecs := os.Getenv("GOST_MQTT_KEEPALIVE_SECS")
	if keepaliveSecs != "" {
		keepalive, err := strconv.Atoi(keepaliveSecs)
		if err == nil {
			conf.MQTT.KeepAliveSec = keepalive
		}
	}

	pingTimeoutSecs := os.Getenv("GOST_MQTT_PINGTIMEOUT_SECS")
	if pingTimeoutSecs != "" {
		pingTimeout, err := strconv.Atoi(pingTimeoutSecs)
		if err == nil {
			conf.MQTT.PingTimeoutSec = pingTimeout
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
