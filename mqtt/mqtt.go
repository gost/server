package mqtt

import (
	"fmt"
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/gost/server/configuration"
	"github.com/gost/server/sensorthings/models"
)

// MQTT is the implementation of the MQTT server
type MQTT struct {
	host       string
	port       int
	connecting bool
	client     paho.Client
	api        *models.API
}

// CreateMQTTClient creates a new MQTT client
func CreateMQTTClient(config configuration.MQTTConfig) models.MQTTClient {
	opts := paho.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%v", config.Host, config.Port)).SetClientID("gost")
	opts.SetCleanSession(true)
	opts.SetKeepAlive(300 * time.Second)
	opts.SetPingTimeout(20 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(connectionLostHandler)
	opts.SetOnConnectHandler(connectHandler)
	client := paho.NewClient(opts)

	return &MQTT{
		host:   config.Host,
		port:   config.Port,
		client: client,
	}
}

// Start running the MQTT client
func (m *MQTT) Start(api *models.API) {
	m.api = api
	log.Printf("Starting MQTT client on %s", fmt.Sprintf("tcp://%s:%v", m.host, m.port))
	m.connect()

	a := *m.api
	topics := *a.GetTopics()
	for _, t := range topics {
		topic := t
		if token := m.client.Subscribe(topic.Path, 0, func(client paho.Client, msg paho.Message) { go topic.Handler(m.api, msg.Topic(), msg.Payload()) }); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}
	}

	/*
		if token := m.client.Subscribe("Datastreams(1)/Observations", 0, observationHandler); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}
	*/
}

// Stop the MQTT client
func (m *MQTT) Stop() {
	m.client.Disconnect(500)
}

// Publish a message on a topic
func (m *MQTT) Publish(topic string, message string, qos byte) {
	token := m.client.Publish(topic, qos, false, message)
	token.Wait()
}

func (m *MQTT) connect() {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		if !m.connecting {
			log.Printf("MQTT client %v", token.Error())
			m.retryConnect()
		}
	}
}

// retryConnect starts a ticker which tries to connect every xx seconds and stops the ticker
// when a connection is established. This is useful when MQTT Broker and GOST are hosted on the same
// machine and GOST is started before mosquito
func (m *MQTT) retryConnect() {
	log.Printf("MQTT client starting reconnect procedure in background")

	m.connecting = true
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			m.connect()
			if m.client.IsConnected() {
				ticker.Stop()
				m.connecting = false
			}
		}
	}()
}

func connectHandler(c paho.Client) {
	log.Printf("MQTT client connected")
}

//ToDo: bubble up and call retryConnect?
func connectionLostHandler(c paho.Client, err error) {
	log.Printf("MQTT client lost connection: %v", err)
}
