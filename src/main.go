package main

import (
	"flag"
	"log"

	"github.com/geodan/gost/src/configuration"
	"github.com/geodan/gost/src/database/postgis"
	"github.com/geodan/gost/src/http"
	"github.com/geodan/gost/src/mqtt"
	"github.com/geodan/gost/src/sensorthings/api"
	"github.com/geodan/gost/src/sensorthings/models"
)

func main() {
	log.Println("Starting GOST....")
	cfgFlag := flag.String("config", "config.yaml", "path of the config file")
	installFlag := flag.String("install", "", "path to the database creation file")
	flag.Parse()

	cfg := *cfgFlag
	conf, err := configuration.GetConfig(cfg)
	if err != nil {
		log.Fatal("config read error: ", err)
		return
	}

	configuration.SetEnvironmentVariables(&conf)

	database := postgis.NewDatabase(conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Database, conf.Database.Schema, conf.Database.SSL, conf.Database.MaxIdleConns, conf.Database.MaxOpenConns)
	database.Start()

	// if install is supplied create database and close, if not start server
	sqlFile := *installFlag
	if len(sqlFile) != 0 {
		createDatabase(database, sqlFile)
	} else {
		mqttClient := mqtt.CreateMQTTClient(conf.MQTT)
		stAPI := api.NewAPI(database, conf, mqttClient)
		mqttClient.Start(&stAPI)
		createAndStartServer(&stAPI)
	}
}

func createDatabase(db models.Database, sqlFile string) {
	log.Println("--CREATING DATABASE--")

	err := db.CreateSchema(sqlFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database created successfully, you can start your server now")
}

// createAndStartServer creates the GOST HTTPServer and starts it
func createAndStartServer(api *models.API) {
	a := *api
	a.Start()

	gostServer := http.CreateServer(a.GetConfig().Server.Host, a.GetConfig().Server.Port, api)
	gostServer.Start()
}
