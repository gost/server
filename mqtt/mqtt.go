package mqtt

import (
	"os"
	"os/signal"

	"fmt"
	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/surge/glog"
	"github.com/surgemq/surgemq/service"
)

// MQTT is the implementation of the MQTT server
type MQTT struct {
	server           service.Server
	keepAlive        int
	connectTimeout   int
	ackTimeout       int
	timeoutRetries   int
	authenticator    string
	sessionsProvider string
	topicsProvider   string
	cpuprofile       string
	mqttaddr         string
	wsAddr           string // HTTPS websocket address eg. :8080
	wssAddr          string // HTTPS websocket address, eg. :8081
	wssCertPath      string // path to HTTPS public key
	wssKeyPath       string // path to HTTPS private key
}

// NewMQTTServer creates a new MQTT server
func NewMQTTServer(config configuration.MQTTConfig) models.MQTTServer {
	return &MQTT{
		keepAlive:        config.KeepAlive,
		connectTimeout:   config.ConnectTimeout,
		ackTimeout:       config.AckTimeout,
		timeoutRetries:   config.TimeoutRetries,
		authenticator:    service.DefaultAuthenticator,
		sessionsProvider: service.DefaultSessionsProvider,
		topicsProvider:   service.DefaultTopicsProvider,
		mqttaddr:         fmt.Sprintf("tcp://:%v", config.Port),
		cpuprofile:       "",
		wsAddr:           "", // HTTPS websocket address eg. :8080
		wssAddr:          "", // HTTPS websocket address, eg. :8081
		wssCertPath:      "", // path to HTTPS public key
		wssKeyPath:       "", // path to HTTPS private key
	}
}

// Start running the MQTT server
func (m *MQTT) Start() {
	var err error

	m.server = service.Server{
		KeepAlive:        m.keepAlive,
		ConnectTimeout:   m.connectTimeout,
		AckTimeout:       m.ackTimeout,
		TimeoutRetries:   m.timeoutRetries,
		SessionsProvider: m.sessionsProvider,
		TopicsProvider:   m.topicsProvider,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	go func() {
		sig := <-sigChan
		glog.Errorf("Existing due to trapped signal; %v", sig)
		m.server.Close()
		os.Exit(0)
	}()

	if len(m.wsAddr) > 0 || len(m.wssAddr) > 0 {
		/*addr := "tcp://127.0.0.1:1883"
		AddWebsocketHandler("/mqtt", addr)

		if len(wsAddr) > 0 {
			go ListenAndServeWebsocket(wsAddr)
		}

		if len(wssAddr) > 0 && len(wssCertPath) > 0 && len(wssKeyPath) > 0 {
			go ListenAndServeWebsocketSecure(wssAddr, wssCertPath, wssKeyPath)
		}*/
	}

	go func() {
		err = m.server.ListenAndServe(m.mqttaddr)
		if err != nil {
			glog.Errorf("surgemq/main: %v", err)
		}
	}()
}

// Stop the MQTT server
func (m *MQTT) Stop() {
	m.server.Close()
}
