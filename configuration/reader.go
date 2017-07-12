package configuration

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// readFile reads the bytes from a given file
func readFile(cfgFile string) ([]byte, error) {
	source, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	return source, nil
}

// readConfig tries to parse the byte data into a config file
func readConfig(fileContent []byte) (Config, error) {
	config := Config{}
	err := yaml.Unmarshal(fileContent, &config)
	return config, err
}

// GetConfig retrieves a new configuration from the given config file
// Fatal when config does not exist or cannot be read
func GetConfig(cfgFile string) (Config, error) {
	content, err := readFile(cfgFile)
	if err != nil {
		return Config{}, err
	}

	conf, err := readConfig(content)
	if err != nil {
		return Config{}, err
	}

	CurrentConfig = conf
	return CurrentConfig, nil
}
