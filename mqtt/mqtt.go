package mqtt

import (
	"fmt"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/gost/server/configuration"
	gostLog "github.com/gost/server/log"
	"github.com/gost/server/sensorthings/models"
	log "github.com/sirupsen/logrus"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

var logger *log.Entry

// MQTT is the implementation of the MQTT server
type MQTT struct {
	host            string
	port            int
	prefix          string
	clientID        string
	username        string
	password        string
	caCertPath      string
	clientCertPath  string
	privateKeyPath  string
	subscriptionQos byte
	persistent      bool
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

func (m *MQTT) getProtocol() string{
	if(m.clientCertPath != "" && m.privateKeyPath != ""){
		return "ssl"
	}else {
		return "tcp"
	}
}

func initMQTTClientOptions(client *MQTT) (*paho.ClientOptions, error) {

	opts := paho.NewClientOptions() // uses defaults: https://godoc.org/github.com/eclipse/paho.mqtt.golang#NewClientOptions

	if client.username != "" {
		opts.SetUsername(client.username)
	}
	if client.password != "" {
		opts.SetPassword(client.password)
	}

	// TLS CONFIG
	tlsConfig := &tls.Config{}
	if client.caCertPath != "" {

		// Import trusted certificates from CAfile.pem.
		// Alternatively, manually add CA certificates to
		// default openssl CA bundle.
		tlsConfig.RootCAs = x509.NewCertPool()
		pemCerts, err := ioutil.ReadFile(client.caCertPath)
		if err == nil {
			tlsConfig.RootCAs.AppendCertsFromPEM(pemCerts)
		}
	}
	if client.clientCertPath != "" && client.privateKeyPath != "" {
		// Import client certificate/key pair
		cert, err := tls.LoadX509KeyPair(client.clientCertPath, client.privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("error loading client keypair: %s", err)
		}
		// Just to print out the client certificate..
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing client certificate: %s", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	opts.AddBroker(fmt.Sprintf("%s://%s:%v",client.getProtocol() , client.host, client.port))
	opts.SetTLSConfig(tlsConfig)

	opts.SetClientID(client.clientID)
	opts.SetCleanSession(!client.persistent)
	opts.SetKeepAlive(300 * time.Second)
	opts.SetPingTimeout(20 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(client.connectionLostHandler)
	opts.SetOnConnectHandler(client.connectHandler)
	return opts, nil
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
		persistent:      config.Persistent,
		username:        config.Username,
		password:        config.Password,
		caCertPath:      config.CaCertPath,
		clientCertPath:  config.ClientCertPath,
		privateKeyPath:  config.PrivateKeyPath,
	}

	opts,err := initMQTTClientOptions(mqttClient)
	if err != nil {
		logger.Errorf("unable to configure MQTT client: %s", err)
	}

	pahoClient := paho.NewClient(opts)
	mqttClient.client = pahoClient

	return mqttClient
}

// Start running the MQTT client
func (m *MQTT) Start(api *models.API) {
	m.api = api
	logger.Infof("Starting MQTT client on %s", fmt.Sprintf("%s://%s:%v", m.getProtocol(),m.host, m.port))
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
	if !m.disconnected || (m.disconnected && !m.persistent) {
		m.subscribe()
	}

	m.disconnected = false
}

//ToDo: bubble up and call retryConnect?
func (m *MQTT) connectionLostHandler(c paho.Client, err error) {
	logger.Warnf("MQTT client lost connection: %v", err)
	m.disconnected = true
}
