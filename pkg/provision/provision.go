package provision

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofrs/uuid"
	"github.com/mydevicesiot/dps-client/internal/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const digiCertBaltimoreRootCA = `
-----BEGIN CERTIFICATE-----
MIIDdzCCAl+gAwIBAgIEAgAAuTANBgkqhkiG9w0BAQUFADBaMQswCQYDVQQGEwJJ
RTESMBAGA1UEChMJQmFsdGltb3JlMRMwEQYDVQQLEwpDeWJlclRydXN0MSIwIAYD
VQQDExlCYWx0aW1vcmUgQ3liZXJUcnVzdCBSb290MB4XDTAwMDUxMjE4NDYwMFoX
DTI1MDUxMjIzNTkwMFowWjELMAkGA1UEBhMCSUUxEjAQBgNVBAoTCUJhbHRpbW9y
ZTETMBEGA1UECxMKQ3liZXJUcnVzdDEiMCAGA1UEAxMZQmFsdGltb3JlIEN5YmVy
VHJ1c3QgUm9vdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKMEuyKr
mD1X6CZymrV51Cni4eiVgLGw41uOKymaZN+hXe2wCQVt2yguzmKiYv60iNoS6zjr
IZ3AQSsBUnuId9Mcj8e6uYi1agnnc+gRQKfRzMpijS3ljwumUNKoUMMo6vWrJYeK
mpYcqWe4PwzV9/lSEy/CG9VwcPCPwBLKBsua4dnKM3p31vjsufFoREJIE9LAwqSu
XmD+tqYF/LTdB1kC1FkYmGP1pWPgkAx9XbIGevOF6uvUA65ehD5f/xXtabz5OTZy
dc93Uk3zyZAsuT3lySNTPx8kmCFcB5kpvcY67Oduhjprl3RjM71oGDHweI12v/ye
jl0qhqdNkNwnGjkCAwEAAaNFMEMwHQYDVR0OBBYEFOWdWTCCR1jMrPoIVDaGezq1
BE3wMBIGA1UdEwEB/wQIMAYBAf8CAQMwDgYDVR0PAQH/BAQDAgEGMA0GCSqGSIb3
DQEBBQUAA4IBAQCFDF2O5G9RaEIFoN27TyclhAO992T9Ldcw46QQF+vaKSm2eT92
9hkTI7gQCvlYpNRhcL0EYWoSihfVCr3FvDB81ukMJY2GQE/szKN+OMY3EU/t3Wgx
jkzSswF07r51XgdIGn9w/xZchMB5hbgF/X++ZRGjD8ACtPhSNzkE1akxehi/oCr0
Epn3o0WC4zxe9Z2etciefC7IpJ5OCBRLbf1wbWsaY71k5h+3zvDyny67G7fyUIhz
ksLi4xaNmjICq44Y3ekQEe5+NauQrz4wlHrQMz2nZQ/1/I6eYs9HRCwBXbsdtTLS
R9I4LtD+gdwyah617jzV/OeBHRnDJELqYzmp
-----END CERTIFICATE-----
`

// Client connects to the provisioning server and provisions devices.
type Client struct {
	client                MQTT.Client
	opts                  Options
	messageChan           chan message
	registrationStateChan chan map[string]interface{}
	requestScheduled      bool
}

// Options contains the device info used for provisioning.
type Options struct {
	Endpoint       string
	Scope          string
	RegistrationID string
	Cert           string
	Key            string
	OutputFile     string
}

type message struct {
	params     url.Values
	payload    map[string]interface{}
	statusCode int
}

// NewClient creates a new provisioning client.
func NewClient(opts Options) *Client {
	if opts.Endpoint == "" || opts.Scope == "" || opts.RegistrationID == "" || opts.Cert == "" || opts.Key == "" || opts.OutputFile == "" {
		log.WithFields(log.Fields{
			"endpoint":       opts.Endpoint,
			"scope":          opts.Scope,
			"registrationID": opts.RegistrationID,
			"cert":           opts.Cert,
			"key":            opts.Key,
			"outputFile":     opts.OutputFile,
		}).Fatal("missing option")
	}
	log.Info("creating client")
	if log.GetLevel() == log.TraceLevel {
		MQTT.DEBUG = logger.New(log.DebugLevel)
		MQTT.WARN = logger.New(log.WarnLevel)
	}
	MQTT.CRITICAL = logger.New(log.ErrorLevel)
	MQTT.ERROR = logger.New(log.ErrorLevel)
	c := Client{
		opts:                  opts,
		messageChan:           make(chan message),
		registrationStateChan: make(chan map[string]interface{}),
		requestScheduled:      false,
	}
	tlsconfig := c.newTLSConfig()
	mqttOpts := MQTT.NewClientOptions()
	server := fmt.Sprintf("ssl://%s:8883", c.opts.Endpoint)
	username := fmt.Sprintf("%s/registrations/%s/api-version=2019-03-31", c.opts.Scope, c.opts.RegistrationID)
	mqttOpts.AddBroker(server)
	mqttOpts.SetClientID(c.opts.RegistrationID)
	mqttOpts.SetUsername(username)
	mqttOpts.SetTLSConfig(tlsconfig)
	mqttOpts.SetConnectRetry(true)
	mqttOpts.SetAutoReconnect(true)
	mqttOpts.SetOnConnectHandler(func(client MQTT.Client) {
		log.Trace("connected")
		go client.Subscribe("$dps/registrations/res/#", 1, nil)
		go c.sendRegisterRequest(0)
	})
	mqttOpts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		log.Debug("connection lost")
	})
	mqttOpts.SetReconnectingHandler(func(client MQTT.Client, opts *MQTT.ClientOptions) {
		log.Debug("reconnecting")
	})
	mqttOpts.SetDefaultPublishHandler(c.messageHandler)

	c.client = MQTT.NewClient(mqttOpts)
	return &c
}

// ProvisionDevice connects to the server and provisions the device.
func (c *Client) ProvisionDevice() {
	err := c.connect()
	if err != nil {
		log.Fatal("connection error")
	}

	go c.messageLoop()
	registrationState := <-c.registrationStateChan

	log.WithFields(log.Fields{
		"hub":       registrationState["assignedHub"],
		"device id": registrationState["deviceId"],
	}).Info("device registered")

	b, err := ioutil.ReadFile(c.opts.OutputFile)
	if err == nil {
		viper.SetConfigType("toml")
		viper.MergeConfig(bytes.NewBuffer(b))
	}
	viper.Set("integration.mqtt.auth.azure_iot_hub.hostname", registrationState["assignedHub"])
	viper.Set("integration.mqtt.auth.azure_iot_hub.device_id", registrationState["deviceId"])
	viper.Set("integration.mqtt.auth.azure_iot_hub.tls_cert", c.opts.Cert)
	viper.Set("integration.mqtt.auth.azure_iot_hub.tls_key", c.opts.Key)
	viper.Set("integration.mqtt.auth.azure_iot_hub.provisioning.Endpoint", c.opts.Endpoint)
	viper.Set("integration.mqtt.auth.azure_iot_hub.provisioning.scope", c.opts.Scope)

	//Create the output file if it doesn't exist
	file, _ := os.OpenFile(c.opts.OutputFile, os.O_RDONLY|os.O_CREATE, 0644)
	file.Close()
	err = viper.WriteConfigAs(c.opts.OutputFile)
	if err != nil {
		log.WithError(err).Error("error writing opts file")
	}

	log.Info("disconnecting")
	c.client.Disconnect(250)
}

func (c *Client) newTLSConfig() *tls.Config {
	certpool := x509.NewCertPool()
	if !certpool.AppendCertsFromPEM([]byte(digiCertBaltimoreRootCA)) {
		log.Fatal("append ca cert from pem error")
	}

	cert, err := tls.LoadX509KeyPair(c.opts.Cert, c.opts.Key)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		RootCAs:      certpool,
		Certificates: []tls.Certificate{cert},
	}
}

func (c *Client) connect() error {
	log.Info("connecting")
	token := c.client.Connect()
	token.Wait()
	return token.Error()
}

func (c *Client) messageLoop() {
	delay := int64(10)
	retryAfter := int64(2)
	for true {
		select {
		case msg := <-c.messageChan:
			if _, hasParam := msg.params["retry-after"]; hasParam {
				retryAfter, _ = strconv.ParseInt(msg.params["retry-after"][0], 10, 64)
			}
			switch {
			case msg.statusCode >= 300:
				log.WithFields(log.Fields{
					"statusCode": msg.statusCode,
				}).Error("incoming message failure")
				if msg.statusCode <= 429 {
					const maxDelay = 1800
					delay += 10
					if delay > maxDelay {
						delay = maxDelay
					}
					retryAfter = delay
				}
				go func(msg message) {
					log.Infof("retry register after %v seconds", retryAfter)
					go c.sendRegisterRequest(time.Duration(retryAfter) * time.Second)
				}(msg)
			default:
				switch msg.payload["status"] {
				case "assigning":
					go func(msg message) {
						time.Sleep(time.Duration(retryAfter) * time.Second)
						rid := uuid.Must(uuid.NewV4())
						topic := fmt.Sprintf("$dps/registrations/GET/iotdps-get-operationstatus/?$rid=%s&operationId=%s", rid, msg.payload["operationId"])
						log.WithFields(log.Fields{
							"topic": topic,
						}).Info("sending operation status message")
						c.client.Publish(topic, 1, false, " ")
					}(msg)
				case "assigned":
					c.registrationStateChan <- msg.payload["registrationState"].(map[string]interface{})
					return
				}
			}
		case <-time.After((time.Duration(retryAfter) + 30) * time.Second):
			log.Error("timed out, retrying request")
			go c.sendRegisterRequest(0)
		}
	}
}

func (c *Client) sendRegisterRequest(delay time.Duration) {
	if c.requestScheduled {
		return
	}
	c.requestScheduled = true
	if delay > 0 {
		time.Sleep(delay)
	}
	rid := uuid.Must(uuid.NewV4())
	topic := fmt.Sprintf("$dps/registrations/PUT/iotdps-register/?$rid=%s", rid)
	payload := fmt.Sprintf(`{"RegistrationId":"%s"}`, c.opts.RegistrationID)
	log.WithFields(log.Fields{
		"topic":   topic,
		"payload": payload,
	}).Info("sending register message")
	c.client.Publish(topic, 1, false, payload)
	c.requestScheduled = false
}

func (c *Client) messageHandler(client MQTT.Client, msg MQTT.Message) {
	const processingError = 600
	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
	}).Info("received message")
	parts := strings.SplitN(msg.Topic(), "/", 5)
	statusCode, err := strconv.Atoi(parts[3])
	if err != nil {
		statusCode = processingError
	}
	log.WithFields(log.Fields{
		"statusCode": statusCode,
	}).Debug("incoming message status")
	paramString := strings.SplitN(parts[4], "$", -1)
	params, err := url.ParseQuery(paramString[1])
	if err != nil {
		statusCode = processingError
	}
	receivedMsg := message{
		params: params,
	}
	if err := json.Unmarshal(msg.Payload(), &receivedMsg.payload); err != nil {
		statusCode = processingError
	}
	receivedMsg.statusCode = statusCode
	c.messageChan <- receivedMsg
}
