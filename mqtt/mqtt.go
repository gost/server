package mqtt

import (
	"fmt"
	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/sensorthings/models"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"log"
)

// MQTT is the implementation of the MQTT server
type MQTT struct {
	host       string
	port       int
	client     paho.Client
	connecting bool
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
func (m *MQTT) Start() {
	log.Printf("Starting MQTT client on %s", fmt.Sprintf("tcp://%s:%v", m.host, m.port))
	m.connect()
	if token := m.client.Subscribe("/gost/#", 0, brokerLoadHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

// Stop the MQTT client
func (m *MQTT) Stop() {
	m.client.Disconnect(500)
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

func connectionLostHandler(c paho.Client, err error) {
	log.Printf("MQTT client lost connection: %v", err)
}

func brokerLoadHandler(client paho.Client, msg paho.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
