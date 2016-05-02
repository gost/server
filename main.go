package main

import (
	"flag"
	"log"

	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/database"
	"github.com/geodan/gost/http"
	//"github.com/geodan/gost/mqtt"
	"github.com/geodan/gost/sensorthings/api"
	"github.com/geodan/gost/sensorthings/models"
)

var (
	stAPI models.API
	//mqttServer mqtt.MQTTServer
)

func main() {
	createAndStartServer(&stAPI)
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

	database := database.NewDatabase(conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Database, conf.Database.Schema, conf.Database.SSL)
	database.Start()

	//mqttServer = mqtt.NewMQTTServer()
	//mqttServer.Start()

	stAPI = api.NewAPI(database, conf)
}

// createAndStartServer creates the GOST HTTPServer and starts it
func createAndStartServer(api *models.API) {
	a := *api
	gostServer := http.NewServer(a.GetConfig().Server.Host, a.GetConfig().Server.Port, api)
	gostServer.Start()
}
