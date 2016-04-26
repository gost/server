package configuration

import (
	"os"
	"encoding/json"
	"log"
	"errors"
)

// GetConfig retrieves a new configuration from the given config file
// Fatal when config does not exist or cannot be read
func GetConfig(cfgFile string) Config {
	conf, confError := readConfig(cfgFile);
	if confError != nil {
		log.Fatal("config read error: ", confError)
	}

	return conf;
}

func readConfig(cfgFile string) (Config, error) {
	var configuration Config

	if _, err := os.Stat(cfgFile);
	os.IsNotExist(err) {
		return configuration, errors.New("file not found")
	}

	file, _ := os.Open(cfgFile)
	decoder := json.NewDecoder(file)
	configuration = Config{}
	err := decoder.Decode(&configuration)
	file.Close()

	return configuration, err;
}