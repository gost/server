package http

import (
	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/database/postgis"
	"github.com/geodan/gost/mqtt"
	"github.com/geodan/gost/sensorthings/api"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateServer(t *testing.T) {
	// arrange
	server := createTestServer(8080, false)

	// act
	server.Stop()

	// assert
	assert.NotNil(t, server)
}

func TestFailRunServerHttp(t *testing.T) {
	// arrange
	server := createTestServer(789456456, false)

	//assert
	assert.Panics(t, func() { server.Start() })
}

func TestFailRunServerHttps(t *testing.T) {
	// arrange
	server := createTestServer(8080, true)

	//assert
	assert.Panics(t, func() { server.Start() })
}

func TestLowerCaseURI(t *testing.T) {
	n := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.True(t, req.URL.Path == "/test")
	})
	ts := httptest.NewServer(LowerCaseURI(n))
	defer ts.Close()
	res, err := http.Get(ts.URL + "/TEST")
	if err == nil && res != nil {
		defer res.Body.Close()
		b, _ := ioutil.ReadAll(res.Body)
		assert.NotNil(t, b)
	} else {
		t.Fail()
	}
}

func TestPostProcessHandler(t *testing.T) {
	n := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusTeapot)
		rw.Header().Add("Location", "tea location")
		rw.Write([]byte("hello teapot"))
	})
	ts := httptest.NewServer(PostProcessHandler(n, "localhost"))
	defer ts.Close()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("X-Forwarded-For", "coffee")
	res, err := client.Do(req)
	if err == nil && res != nil {
		defer res.Body.Close()
		b, _ := ioutil.ReadAll(res.Body)
		body := string(b)
		assert.NotNil(t, body)
		assert.True(t, body == "hello teapot")
		assert.True(t, res.StatusCode == http.StatusTeapot)
		assert.True(t, res.Header.Get("Location") == "tea location")
	} else {
		t.Fail()
	}
}

func createTestServer(port int, https bool) Server {
	cfg := configuration.Config{
		Server: configuration.ServerConfig{ExternalURI: "http://localhost:8080/"},
	}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100, 200)
	stAPI := api.NewAPI(database, cfg, mqttServer)
	return CreateServer("localhost", port, &stAPI, https, "", "")
}
