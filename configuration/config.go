package configuration

import (
	"fmt"
	"strings"
)

// Config holds the settings for the Http server and databases
type Config struct {
	Server struct {
		Name        string `yaml:"name"`
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		ExternalUri string `yaml:"externalUri"`
	}
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Schema   string `yaml:"schema"`
		SSL      bool   `yaml:"ssl"`
	}
}

// GetInternalServerUri gets the internal Http server address
// for example: "localhost:8080"
func (c *Config) GetInternalServerUri() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetExternalServerUri gets the external Http server address, trailing slash is removed when present in Config.Server.ExternalUri
// for example "http://www.mysensorplatform"
func (c *Config) GetExternalServerUri() string {
	return fmt.Sprintf("%s", strings.Trim(c.Server.ExternalUri, "/"))
}
