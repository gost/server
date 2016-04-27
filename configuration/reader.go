package configuration

import (
	"os"
	"log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func readConfig(cfgFile string) (Config, error) {
	config := Config{}

	if _, err := os.Stat(cfgFile);
	os.IsNotExist(err) {
		return config, err
	}

	source, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return config, err
	}

	err2 := yaml.Unmarshal(source, &config)
	if err2 != nil {
		return config, err2
	}

	return config, nil
}

// GetConfig retrieves a new configuration from the given config file
// Fatal when config does not exist or cannot be read
func GetConfig(cfgFile string) Config {
	conf, confError := readConfig(cfgFile);
	if confError != nil {
		log.Fatal("config read error: ", confError)
	}

	return conf;
}