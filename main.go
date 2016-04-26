package main

import (
	"flag"

	"github.com/tebben/gost/configuration"
	"github.com/tebben/gost/gostdb"
	"github.com/tebben/gost/gosthttp"
	//"github.com/tebben/gost/mqtt"
	"github.com/tebben/gost/sensorthings"
)

var (
	conf       configuration.Config
	api        sensorthings.SensorThingsApi
	//mqttServer mqtt.MQTTServer
)

func main() {
	createAndStartServer(&api)
}

// Initialises GOST
// Parse flags, read config, setup database and api
func init() {
	var cfgFlag = "config.json"
	flag.StringVar(&cfgFlag, "config", "config.json", "path of the config file, default = config.json")
	flag.Parse()

	conf = configuration.GetConfig(cfgFlag)
	database := gostdb.NewDatabase(conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Database, conf.Database.Schema, conf.Database.SSL)
	database.Start()

	//mqttServer = mqtt.NewMQTTServer()
	//mqttServer.Start()

	api = sensorthings.NewApi(database, conf)
}

// createAndStartServer creates the GOST HTTPServer and starts it
func createAndStartServer(api *sensorthings.SensorThingsApi) {
	gostServer := gosthttp.NewServer(conf.Server.Host, conf.Server.Port, api)
	gostServer.Start()
}
