package configuration

import (
	"fmt"
	"strings"
)

// CurrentConfig will be set after loading so it can be accessed from outside
var CurrentConfig Config

// Config contains the settings for the Http server, databases and mqtt
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	MQTT     MQTTConfig     `yaml:"mqtt"`
	Logger   LoggerConfig   `yaml:"logger"`
}

// ServerConfig contains the general server information
type ServerConfig struct {
	Name              string `yaml:"name"`
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	ExternalURI       string `yaml:"externalUri"`
	HTTPS             bool   `yaml:"https"`
	HTTPSCert         string `yaml:"httpsCert"`
	HTTPSKey          string `yaml:"httpsKey"`
	MaxEntityResponse int    `yaml:"maxEntityResponse"`
	IndentedJSON      bool   `yaml:"indentedJson"`
}

// DatabaseConfig contains the database server information, can be overruled by environment variables
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	Schema       string `yaml:"schema"`
	SSL          bool   `yaml:"ssl"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
}

// MQTTConfig contains the MQTT client information
type MQTTConfig struct {
	Enabled         bool   `yaml:"enabled"`
	Verbose         bool   `yaml:"verbose"`
	Host            string `yaml:"host"`
	Prefix          string `yaml:"prefix"`
	ClientID        string `yaml:"clientId"`
	Port            int    `yaml:"port"`
	SubscriptionQos byte   `yaml:"subscriptionQos"`
	Persistent      bool   `yaml:"persistent"`
	Order           bool   `yaml:"order"`
	SSL             bool   `yaml:"ssl"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	CaCertFile      string `yaml:"caCertFile"`
	ClientCertFile  string `yaml:"clientCertFile"`
	PrivateKeyFile  string `yaml:"privateKeyFile"`
	KeepAliveSec	int    `yaml:"keepAliveSec"`
	PingTimeoutSec	int    `yaml:"pingTimeoutSec"`
}

// LoggerConfig contains the logging configuration used to initialize the logger
type LoggerConfig struct {
	FileName string `yaml:"fileName"`
	Verbose  bool   `yaml:"verbose"`
}

// GetInternalServerURI gets the internal Http server address
// for example: "localhost:8080"
func (c *Config) GetInternalServerURI() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetExternalServerURI gets the external Http server address, trailing slash is removed when present in Config.Server.ExternalUri
// for example "http://www.mysensorplatform"
func (c *Config) GetExternalServerURI() string {
	return strings.Trim(c.Server.ExternalURI, "/")
}
