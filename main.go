package main

import (
	"flag"

	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/gostdb"
	"github.com/geodan/gost/gosthttp"
	//"github.com/geodan/gost/mqtt"
	"log"

	"github.com/geodan/gost/sensorthings"
)

var (
	api sensorthings.SensorThingsAPI
	//mqttServer mqtt.MQTTServer
)

func main() {
	createAndStartServer(&api)
}

// Initialises GOST
// Parse flags, read config, setup database and api
func init() {
	var cfgFlag = "config.yaml"
	flag.StringVar(&cfgFlag, "config", "config.yaml", "path of the config file, default = config.yaml")
	flag.Parse()

	conf, err := configuration.GetConfig(cfgFlag)
	if err != nil {
		log.Fatal("config read error: ", err)
	}

	database := gostdb.NewDatabase(conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Database, conf.Database.Schema, conf.Database.SSL)
	database.Start()

	//mqttServer = mqtt.NewMQTTServer()
	//mqttServer.Start()

	api = sensorthings.NewAPI(database, conf)
}

// createAndStartServer creates the GOST HTTPServer and starts it
func createAndStartServer(api *sensorthings.SensorThingsAPI) {
	a := *api
	gostServer := gosthttp.NewServer(a.GetConfig().Server.Host, a.GetConfig().Server.Port, api)
	gostServer.Start()
}
