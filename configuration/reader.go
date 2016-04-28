package configuration

import (
	"os"
	"log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func readFile(cfgFile string) ([]byte, error){
	if _, err := os.Stat(cfgFile);
	os.IsNotExist(err) {
		return nil, err
	}

	source, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
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
func GetConfig(cfgFile string) Config {
	content, err := readFile(cfgFile)
	if err != nil {
		log.Fatal("config read error: ", err)
	}

	conf, err2 := readConfig(content);
	if err2 != nil {
		log.Fatal("config read error: ", err2)
	}

	return conf;
}