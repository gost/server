package http

import (
	"github.com/gorilla/mux"
	"github.com/gost/server/configuration"
	"github.com/gost/server/database/postgis"
	"github.com/gost/server/mqtt"
	"github.com/gost/server/sensorthings/api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router *mux.Router
var req *http.Request
var respRec *httptest.ResponseRecorder

func setup() {
	// arrange
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100, 200)
	a := api.NewAPI(database, cfg, mqttServer)
	router = CreateRouter(&a)

	//The response recorder used to record HTTP responses
	respRec = httptest.NewRecorder()
}

// Test the router functionality
func TestCreateRouter(t *testing.T) {
	// arrange
	setup()

	// assert
	assert.NotNil(t, router, "Router should be created")
}

func TestEndpoints(t *testing.T) {
	req, _ = http.NewRequest("GET", "/v1.0", nil)
	router.ServeHTTP(respRec, req)
	if respRec.Code != http.StatusOK {
		t.Fatal("Server endpoint /v1.0 error: Returned ", respRec.Code, " instead of ", http.StatusOK)
	}
}
