package configuration

import (
	"fmt"
	"strings"
)

// GOST configuration, holding the settings for the Http server and databases
type Config struct {
	Server struct {
		Name string `json:"name"`
		Host string `json:"host"`
		Port int `json:"port"`
		ExternalUri string `json:"externalUri"`
	} `json:"server"`
	Database struct {
	       Host string `json:"host"`
	       Port int `json:"port"`
	       User string `json:"user"`
	       Password string `json:"password"`
	       Database string `json:"database"`
	       Schema string `json:"schema"`
	       SSL bool `json:"ssl"`
       } `json:"database"`
}

// Get the internal Http server address
// for example: "localhost:8080"
func (c *Config) GetInternalServerUri() string{
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// Get the external Http server address, trailing slash is removed when present in Config.Server.ExternalUri
// for example "http://www.mysensorplatform"
func (c *Config) GetExternalServerUri() string{
	return fmt.Sprintf("%s", strings.Trim(c.Server.ExternalUri, "/"))
}