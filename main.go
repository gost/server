package main

import (
	"flag"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/gost/server/configuration"
	"github.com/gost/server/database/postgis"
	"github.com/gost/server/http"
	"github.com/gost/server/mqtt"
	"github.com/gost/server/sensorthings/api"
	"github.com/gost/server/sensorthings/models"
)

var (
	stAPI      models.API
	gostServer http.Server
	mqttClient models.MQTTClient
	file       *os.File
	gostLogger = log.New()
)

func main() {
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-stop
		cleanup()
		log.Print("GOST stopped gracefully")
		os.Exit(1)
	}()

	log.Println("Starting GOST....")
	log.Println("Showing debug logs")

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
	err = initializeLogger(file, conf.Logger.FileName, gostLogger, new(log.TextFormatter), conf.Logger.Verbose)
	if err != nil {
		fmt.Println("Error initializing logger, defaulting to stdout. Error: " + err.Error())
	}
	database := postgis.NewDatabase(
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Database,
		conf.Database.Schema,
		conf.Database.SSL,
		conf.Database.MaxIdleConns,
		conf.Database.MaxOpenConns,
		conf.Server.MaxEntityResponse)
	go database.Start()

	// if install is supplied create database and close, if not start server
	sqlFile := *installFlag
	if len(sqlFile) != 0 {
		createDatabase(database, sqlFile)
	} else {
		mqttClient = mqtt.CreateMQTTClient(conf.MQTT)
		stAPI = api.NewAPI(database, conf, mqttClient)

		if conf.MQTT.Enabled {
			mqttClient.Start(&stAPI)
		}

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

	config := a.GetConfig()
	gostServer = http.CreateServer(
		config.Server.Host,
		config.Server.Port,
		api,
		config.Server.HTTPS,
		config.Server.HTTPSCert,
		config.Server.HTTPSKey)
	gostServer.Start()
}

func cleanup() {
	if gostServer != nil {
		gostServer.Stop()
	}

	file.Close()
}

func initializeLogger(file *os.File, logFileName string, logger *log.Logger, format log.Formatter, verboseFlag bool) error {

	file, err := os.OpenFile(logFileName+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Log file cannot be opened")
		file.Close()
		file = os.Stdout
	}

	logger.Out = file
	if format != nil {
		logger.Formatter = format
	}

	if verboseFlag {
		logger.Level = log.DebugLevel
	} else {
		logger.Level = log.WarnLevel
	}

	return err
}
