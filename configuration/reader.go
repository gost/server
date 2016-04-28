package configuration

import (
	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func readFile(cfgFile string) ([]byte, error){
	if _, err := os.Stat(cfgFile);
	os.IsNotExist(err) {
		return nil, err
	}

	source, err2 := ioutil.ReadFile(cfgFile)
	if err2 != nil {
		return nil, err2
	}

	return source, nil
}

func readConfig(fileContent []byte) (Config, error) {
	config := Config{}
	err2 := yaml.Unmarshal(fileContent, &config)
	if err2 != nil {
		return config, err2
	}

	return config, nil
}

// GetConfig retrieves a new configuration from the given config file
// Fatal when config does not exist or cannot be read
func GetConfig(cfgFile string) (Config, error) {
	content, err := readFile(cfgFile)
	if err != nil {
		return Config{}, err
	}

	conf, err := readConfig(content);
	if err != nil {
		return Config{}, err
	}

	return conf, nil;
}