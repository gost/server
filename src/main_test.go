package main

import (
	"testing"

	"github.com/geodan/gost/src/configuration"
	"github.com/geodan/gost/src/database/postgis"
	"github.com/geodan/gost/src/mqtt"
	"github.com/geodan/gost/src/sensorthings/api"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestUtils(t *testing.T) {
	// todo add utils tests...
	// arrange
	var a = 1
	var b = 2
	// act
	var res = a + b
	// assert
	assert.Equal(t, 3, res, "computer error again")
}

var jsonObservationMinimal = `{ "result": 18 }`

var configData = `
server:
    name: GOST Server
    host: localhost
    port: 8080
    externalUri: localhost:8080/
database:
    host: localhost
    port: 5432
    user: postgres
    password: postgres
    database: gost
    schema: v1
    ssl: false
mqtt:
    enabled: true
    host: test.mosquitto.org
    port: 1883
`

var (
	jsonBytes []byte
	testAPI   models.API
)

func init() {
	jsonBytes = []byte(jsonObservationMinimal)

	content := []byte(configData)
	conf := configuration.Config{}
	yaml.Unmarshal(content, &conf)

	database := postgis.NewDatabase(conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Database, conf.Database.Schema, conf.Database.SSL, conf.Database.MaxIdleConns, conf.Database.MaxOpenConns)
	database.Start()

	mqttClient := mqtt.CreateMQTTClient(conf.MQTT)
	testAPI = api.NewAPI(database, conf, mqttClient)
}

func BenchmarkIncomingObservation(b *testing.B) {
	//b.StopTimer()
	//b.StartTimer()

	for i := 0; i < b.N; i++ {
		createAndParseObservation()
	}
}

func createAndParseObservation() {
	o := &entities.Observation{}
	o.ParseEntity(jsonBytes)
	testAPI.PostObservationByDatastream("1", o)
}

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}
