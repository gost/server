package mqtt

import (
	"fmt"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/gost/server/configuration"
	gostLog "github.com/gost/server/log"
	"github.com/gost/server/sensorthings/models"
	log "github.com/sirupsen/logrus"
)

var logger *log.Entry

// MQTT is the implementation of the MQTT server
type MQTT struct {
	host            string
	port            int
	prefix          string
	clientID        string
	subscriptionQos byte
	persitent       bool
	connecting      bool
	disconnected    bool
	client          paho.Client
	api             *models.API
}

func setupLogger() {
	l, err := gostLog.GetLoggerInstance()
	if err != nil {
		log.Error(err)
	}

	logger = l.WithFields(log.Fields{"package": "gost.server.mqtt"})
}

// CreateMQTTClient creates a new MQTT client
func CreateMQTTClient(config configuration.MQTTConfig) models.MQTTClient {
	setupLogger()

	mqttClient := &MQTT{
		host:            config.Host,
		port:            config.Port,
		prefix:          config.Prefix,
		clientID:        config.ClientID,
		subscriptionQos: config.SubscriptionQos,
		persitent:       config.Persistent,
	}

	opts := paho.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%v", config.Host, config.Port)).SetClientID(config.ClientID)
	opts.SetCleanSession(!config.Persistent)
	opts.SetKeepAlive(300 * time.Second)
	opts.SetPingTimeout(20 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(mqttClient.connectionLostHandler)
	opts.SetOnConnectHandler(mqttClient.connectHandler)

	pahoClient := paho.NewClient(opts)
	mqttClient.client = pahoClient

	return mqttClient
}

// Start running the MQTT client
func (m *MQTT) Start(api *models.API) {
	m.api = api
	logger.Infof("Starting MQTT client on %s", fmt.Sprintf("tcp://%s:%v", m.host, m.port))
	m.connect()
}

// Stop the MQTT client
func (m *MQTT) Stop() {
	m.client.Disconnect(500)
}

func (m *MQTT) subscribe() {
	a := *m.api
	topics := *a.GetTopics(m.prefix)

	for _, t := range topics {
		topic := t
		if token := m.client.Subscribe(topic.Path, m.subscriptionQos, func(client paho.Client, msg paho.Message) {
			go topic.Handler(m.api, m.prefix, msg.Topic(), msg.Payload())
		}); token.Wait() && token.Error() != nil {
			logger.Error(token.Error())
		}
	}
}

// Publish a message on a topic
func (m *MQTT) Publish(topic string, message string, qos byte) {
	token := m.client.Publish(topic, qos, false, message)
	token.Wait()
}

func (m *MQTT) connect() {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		if !m.connecting {
			logger.Errorf("MQTT client %v", token.Error())
			m.retryConnect()
		}
	}
}

// retryConnect starts a ticker which tries to connect every xx seconds and stops the ticker
// when a connection is established. This is useful when MQTT Broker and GOST are hosted on the same
// machine and GOST is started before mosquito
func (m *MQTT) retryConnect() {
	logger.Infof("MQTT client starting reconnect procedure in background")
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

func (m *MQTT) connectHandler(c paho.Client) {
	logger.Infof("MQTT client connected")

	// on first connect or connection lost and persistance is off
	if !m.disconnected || (m.disconnected && !m.persitent) {
		m.subscribe()
	}

	m.disconnected = false
}

//ToDo: bubble up and call retryConnect?
func (m *MQTT) connectionLostHandler(c paho.Client, err error) {
	logger.Warnf("MQTT client lost connection: %v", err)
	m.disconnected = true
}
