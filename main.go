package main

import (
	"flag"
	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/database"
	"github.com/geodan/gost/http"
	"log"
	"os"
	"strconv"
	//"github.com/geodan/gost/mqtt"
	"github.com/geodan/gost/sensorthings/api"
	"github.com/geodan/gost/sensorthings/models"
)

var (
	stAPI models.API
	//mqttServer mqtt.MQTTServer
)

func main() {
	cfgFlag := flag.String("config", "config.yaml", "path of the config file")
	installFlag := flag.String("install", "", "path to the database creation file")
	flag.Parse()

	cfg := *cfgFlag
	conf, err := configuration.GetConfig(cfg)
	if err != nil {
		log.Fatal("config read error: ", err)
		return
	}
	
	gostDbHost := os.Getenv("gost_db_host")
	if(gostDbHost != ""){
		conf.Database.Host = gostDbHost;
	}
	gostDbPort := os.Getenv("gost_db_port")
	if(gostDbPort != ""){
		port, err := strconv.Atoi(gostDbPort)
		if(err==nil){
			conf.Database.Port = int(port);
		}
	}
	gostDbUser := os.Getenv("gost_db_user")
	if(gostDbUser != ""){
		conf.Database.User = gostDbUser;
	}
	gostDbPassword := os.Getenv("gost_db_password")
	if(gostDbPassword != ""){
		conf.Database.Password = gostDbPassword;
	}	

	database := database.NewDatabase(conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Database, conf.Database.Schema, conf.Database.SSL)
	database.Start()

	// if install is supplied create database and close, if not start server
	sqlFile := *installFlag
	if len(sqlFile) != 0 {
		createDatabase(database, sqlFile)
	} else {
		//mqttServer = mqtt.NewMQTTServer()
		//mqttServer.Start()
		stAPI = api.NewAPI(database, conf)
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
	gostServer := http.NewServer(a.GetConfig().Server.Host, a.GetConfig().Server.Port, api)
	gostServer.Start()
}
