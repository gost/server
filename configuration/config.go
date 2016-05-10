package configuration

import (
	"fmt"
	"strings"
)

// Config holds the settings for the Http server, databases and mqtt
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	MQTT     MQTTConfig     `yaml:"mqtt"`
}

type ServerConfig struct {
	Name        string `yaml:"name"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	ExternalURI string `yaml:"externalUri"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
	SSL      bool   `yaml:"ssl"`
}

type MQTTConfig struct {
	Enabled bool `yaml:"enabled"`
	host    bool `yaml:"host"`
	Port    int  `yaml:"port"`
}

// GetInternalServerURI gets the internal Http server address
// for example: "localhost:8080"
func (c *Config) GetInternalServerURI() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetExternalServerURI gets the external Http server address, trailing slash is removed when present in Config.Server.ExternalUri
// for example "http://www.mysensorplatform"
func (c *Config) GetExternalServerURI() string {
	return fmt.Sprintf("%s", strings.Trim(c.Server.ExternalURI, "/"))
}
