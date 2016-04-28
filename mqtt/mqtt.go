package mqtt

import (
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/surge/glog"
	"github.com/surgemq/surgemq/service"
)

// Server interface defines the needed MQTT operations
type Server interface {
	Start()
	Stop()
}

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
	wsAddr           string // HTTPS websocket address eg. :8080
	wssAddr          string // HTTPS websocket address, eg. :8081
	wssCertPath      string // path to HTTPS public key
	wssKeyPath       string // path to HTTPS private key
}

// NewMQTTServer creates a new MQTT server
func NewMQTTServer() Server {
	return &MQTT{
		keepAlive:        service.DefaultKeepAlive,
		connectTimeout:   service.DefaultConnectTimeout,
		ackTimeout:       service.DefaultAckTimeout,
		timeoutRetries:   service.DefaultTimeoutRetries,
		authenticator:    service.DefaultAuthenticator,
		sessionsProvider: service.DefaultSessionsProvider,
		topicsProvider:   service.DefaultTopicsProvider,
		cpuprofile:       "",
		wsAddr:           "", // HTTPS websocket address eg. :8080
		wssAddr:          "", // HTTPS websocket address, eg. :8081
		wssCertPath:      "", // path to HTTPS public key
		wssKeyPath:       "", // path to HTTPS private key
	}
}

// Start running the MQTT server
func (m *MQTT) Start() {
	mqttaddr := "tcp://:1883"
	var err error

	//	if m.server == nil {
	m.server = service.Server{
		KeepAlive:        m.keepAlive,
		ConnectTimeout:   m.connectTimeout,
		AckTimeout:       m.ackTimeout,
		TimeoutRetries:   m.timeoutRetries,
		SessionsProvider: m.sessionsProvider,
		TopicsProvider:   m.topicsProvider,
	}

	var f *os.File

	if m.cpuprofile != "" {
		f, err = os.Create(m.cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill)
	go func() {
		sig := <-sigchan
		glog.Errorf("Existing due to trapped signal; %v", sig)

		if f != nil {
			glog.Errorf("Stopping profile")
			pprof.StopCPUProfile()
			f.Close()
		}

		m.server.Close()

		os.Exit(0)
	}()

	if len(m.wsAddr) > 0 || len(m.wssAddr) > 0 {

	}

	go func() {
		err = m.server.ListenAndServe(mqttaddr)
		if err != nil {
			glog.Errorf("surgemq/main: %v", err)
		}
	}()
}

// Stop the MQTT server
func (m *MQTT) Stop() {
	m.server.Close()
}
