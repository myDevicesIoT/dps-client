package provision

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/mydevicesiot/dps-client/internal/logger"
	"github.com/pelletier/go-toml"
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

const microsoftRSARootCA2017 = `
-----BEGIN CERTIFICATE-----
MIIFqDCCA5CgAwIBAgIQHtOXCV/YtLNHcB6qvn9FszANBgkqhkiG9w0BAQwFADBl
MQswCQYDVQQGEwJVUzEeMBwGA1UEChMVTWljcm9zb2Z0IENvcnBvcmF0aW9uMTYw
NAYDVQQDEy1NaWNyb3NvZnQgUlNBIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5
IDIwMTcwHhcNMTkxMjE4MjI1MTIyWhcNNDIwNzE4MjMwMDIzWjBlMQswCQYDVQQG
EwJVUzEeMBwGA1UEChMVTWljcm9zb2Z0IENvcnBvcmF0aW9uMTYwNAYDVQQDEy1N
aWNyb3NvZnQgUlNBIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5IDIwMTcwggIi
MA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDKW76UM4wplZEWCpW9R2LBifOZ
Nt9GkMml7Xhqb0eRaPgnZ1AzHaGm++DlQ6OEAlcBXZxIQIJTELy/xztokLaCLeX0
ZdDMbRnMlfl7rEqUrQ7eS0MdhweSE5CAg2Q1OQT85elss7YfUJQ4ZVBcF0a5toW1
HLUX6NZFndiyJrDKxHBKrmCk3bPZ7Pw71VdyvD/IybLeS2v4I2wDwAW9lcfNcztm
gGTjGqwu+UcF8ga2m3P1eDNbx6H7JyqhtJqRjJHTOoI+dkC0zVJhUXAoP8XFWvLJ
jEm7FFtNyP9nTUwSlq31/niol4fX/V4ggNyhSyL71Imtus5Hl0dVe49FyGcohJUc
aDDv70ngNXtk55iwlNpNhTs+VcQor1fznhPbRiefHqJeRIOkpcrVE7NLP8TjwuaG
YaRSMLl6IE9vDzhTyzMMEyuP1pq9KsgtsRx9S1HKR9FIJ3Jdh+vVReZIZZ2vUpC6
W6IYZVcSn2i51BVrlMRpIpj0M+Dt+VGOQVDJNE92kKz8OMHY4Xu54+OU4UZpyw4K
UGsTuqwPN1q3ErWQgR5WrlcihtnJ0tHXUeOrO8ZV/R4O03QK0dqq6mm4lyiPSMQH
+FJDOvTKVTUssKZqwJz58oHhEmrARdlns87/I6KJClTUFLkqqNfs+avNJVgyeY+Q
W5g5xAgGwax/Dj0ApQIDAQABo1QwUjAOBgNVHQ8BAf8EBAMCAYYwDwYDVR0TAQH/
BAUwAwEB/zAdBgNVHQ4EFgQUCctZf4aycI8awznjwNnpv7tNsiMwEAYJKwYBBAGC
NxUBBAMCAQAwDQYJKoZIhvcNAQEMBQADggIBAKyvPl3CEZaJjqPnktaXFbgToqZC
LgLNFgVZJ8og6Lq46BrsTaiXVq5lQ7GPAJtSzVXNUzltYkyLDVt8LkS/gxCP81OC
gMNPOsduET/m4xaRhPtthH80dK2Jp86519efhGSSvpWhrQlTM93uCupKUY5vVau6
tZRGrox/2KJQJWVggEbbMwSubLWYdFQl3JPk+ONVFT24bcMKpBLBaYVu32TxU5nh
SnUgnZUP5NbcA/FZGOhHibJXWpS2qdgXKxdJ5XbLwVaZOjex/2kskZGT4d9Mozd2
TaGf+G0eHdP67Pv0RR0Tbc/3WeUiJ3IrhvNXuzDtJE3cfVa7o7P4NHmJweDyAmH3
pvwPuxwXC65B2Xy9J6P9LjrRk5Sxcx0ki69bIImtt2dmefU6xqaWM/5TkshGsRGR
xpl/j8nWZjEgQRCHLQzWwa80mMpkg/sTV9HB8Dx6jKXB/ZUhoHHBk2dxEuqPiApp
GWSZI1b7rCoucL5mxAyE7+WL85MB+GqQk2dLsmijtWKP6T+MejteD+eMuMZ87zf9
dOLITzNy4ZQ5bb0Sr74MTnB8G2+NszKTc0QWbej09+CVgI+WXTik9KveCjCHk9hN
AHFiRSdLOkKEW39lt2c0Ui2cFmuqqNh7o0JMcccMyj6D5KbvtwEwXlGjefVwaaZB
RA+GsCyRxj3qrg+E
-----END CERTIFICATE-----
`

const digiCertGlobalRootG2 = `
-----BEGIN CERTIFICATE-----
MIIDjjCCAnagAwIBAgIQAzrx5qcRqaC7KGSxHQn65TANBgkqhkiG9w0BAQsFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH
MjAeFw0xMzA4MDExMjAwMDBaFw0zODAxMTUxMjAwMDBaMGExCzAJBgNVBAYTAlVT
MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdpY2VydC5j
b20xIDAeBgNVBAMTF0RpZ2lDZXJ0IEdsb2JhbCBSb290IEcyMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuzfNNNx7a8myaJCtSnX/RrohCgiN9RlUyfuI
2/Ou8jqJkTx65qsGGmvPrC3oXgkkRLpimn7Wo6h+4FR1IAWsULecYxpsMNzaHxmx
1x7e/dfgy5SDN67sH0NO3Xss0r0upS/kqbitOtSZpLYl6ZtrAGCSYP9PIUkY92eQ
q2EGnI/yuum06ZIya7XzV+hdG82MHauVBJVJ8zUtluNJbd134/tJS7SsVQepj5Wz
tCO7TG1F8PapspUwtP1MVYwnSlcUfIKdzXOS0xZKBgyMUNGPHgm+F6HmIcr9g+UQ
vIOlCsRnKPZzFBQ9RnbDhxSJITRNrw9FDKZJobq7nMWxM4MphQIDAQABo0IwQDAP
BgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB/wQEAwIBhjAdBgNVHQ4EFgQUTiJUIBiV
5uNu5g/6+rkS7QYXjzkwDQYJKoZIhvcNAQELBQADggEBAGBnKJRvDkhj6zHd6mcY
1Yl9PMWLSn/pvtsrF9+wX3N3KjITOYFnQoQj8kVnNeyIv/iPsGEMNKSuIEyExtv4
NeF22d+mQrvHRAiGfzZ0JFrabA0UWTW98kndth/Jsw1HKj2ZL7tcu7XUIOGZX1NG
Fdtom/DzMNU+MeKNhJ7jitralj41E6Vf8PlwUHBHQRFXGU7Aj64GxJUTFy8bJZ91
8rGOmaFvE7FBcf6IKshPECBV1/MUReXgRPTqh5Uykw7+U0b6LJ3/iyK5S9kJRaTe
pLiaWN0bfVKfjllDiIGknibVb63dDcY3fe0Dkhvld1927jyNxF1WW6LZZm6zNTfl
MrY=
-----END CERTIFICATE-----
`

// Client connects to the provisioning server and provisions devices.
type AzureIoTHubProvisioner struct {
	mqttClient            MQTT.Client
	httpClient            http.Client
	opts                  ProvisionerOptions
	messageChan           chan message
	registrationStateChan chan RegistrationState
	requestScheduled      bool
}

// RegistrationRequest represents the body of a DPS registration request
type RegistrationRequest struct {
	RegistrationID string `json:"registrationId"`
}

// RegistrationState represents the body of a DPS registration response containing the current state of a device
type RegistrationState struct {
	RegistrationID         string `json:"registrationId"`
	CreatedDateTimeUtc     string `json:"createdDateTimeUtc"`
	AssignedHub            string `json:"assignedHub"`
	DeviceID               string `json:"deviceId"`
	Status                 string `json:"status"`
	SubStatus              string `json:"substatus"`
	LastUpdatedDateTimeUtc string `json:"lastUpdatedDateTimeUtc"`
	ETag                   string `json:"etag"`
}

// RegistrationResponse represents the body of a DPS registration response
type RegistrationResponse struct {
	OperationID       string            `json:"operationId"`
	Status            string            `json:"status"`
	RegistrationState RegistrationState `json:"registrationState"`
}

type message struct {
	params               url.Values
	registrationResponse RegistrationResponse
	statusCode           int
}

// NewClient creates a new provisioning client.
func NewAzureIoTHubProvisioner(opts ProvisionerOptions) *AzureIoTHubProvisioner {
	if opts.Endpoint == "" || opts.Scope == "" || opts.RegistrationID == "" || opts.Cert == "" || opts.Key == "" || opts.OutputFile == "" {
		log.WithFields(log.Fields{
			"endpoint":       opts.Endpoint,
			"scope":          opts.Scope,
			"registrationID": opts.RegistrationID,
			"cert":           opts.Cert,
			"key":            opts.Key,
			"outputFile":     opts.OutputFile,
			"protocol":       opts.Protocol,
		}).Fatal("missing option")
	}
	log.Info("creating client")
	c := AzureIoTHubProvisioner{
		opts:                  opts,
		messageChan:           make(chan message),
		registrationStateChan: make(chan RegistrationState),
		requestScheduled:      false,
	}
	tlsconfig := c.newTLSConfig()
	if strings.Contains(opts.Protocol, "http") {
		transport := &http.Transport{TLSClientConfig: c.newTLSConfig()}
		c.httpClient = http.Client{Transport: transport}
	} else {
		if log.GetLevel() == log.TraceLevel {
			MQTT.DEBUG = logger.New(log.DebugLevel)
			MQTT.WARN = logger.New(log.WarnLevel)
		}
		MQTT.CRITICAL = logger.New(log.ErrorLevel)
		MQTT.ERROR = logger.New(log.ErrorLevel)
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

		c.mqttClient = MQTT.NewClient(mqttOpts)
	}
	return &c
}

// ProvisionDevice connects to the server and provisions the device.
func (c *AzureIoTHubProvisioner) ProvisionDevice() {
	err := c.connect()
	if err != nil {
		log.Fatal("connection error")
	}

	go c.messageLoop()
	registrationState := <-c.registrationStateChan

	log.WithFields(log.Fields{
		"hub":       registrationState.AssignedHub,
		"device id": registrationState.DeviceID,
	}).Info("device registered")

	c.writeConfigFile(registrationState)

	if c.mqttClient != nil {
		log.Info("disconnecting")
		c.mqttClient.Disconnect(250)
	}
}

func (c *AzureIoTHubProvisioner) newTLSConfig() *tls.Config {
	certpool := x509.NewCertPool()
	rootCAs := fmt.Sprintf("%s%s%s", digiCertBaltimoreRootCA, microsoftRSARootCA2017, digiCertGlobalRootG2)
	if !certpool.AppendCertsFromPEM([]byte(rootCAs)) {
		log.Fatal("append ca certs from pem error")
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

func (c *AzureIoTHubProvisioner) connect() error {
	if c.mqttClient != nil {
		log.Info("connecting")
		token := c.mqttClient.Connect()
		token.Wait()
		return token.Error()
	}
	go c.sendRegisterRequest(0)
	return nil
}

func (c *AzureIoTHubProvisioner) messageLoop() {
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
				switch msg.registrationResponse.Status {
				case "assigning":
					go c.sendOperationStatusRequest(time.Duration(retryAfter)*time.Second, msg)
				case "assigned":
					c.registrationStateChan <- msg.registrationResponse.RegistrationState
					return
				}
			}
		case <-time.After((time.Duration(retryAfter) + 30) * time.Second):
			log.Error("timed out, retrying request")
			go c.sendRegisterRequest(0)
		}
	}
}

func (c *AzureIoTHubProvisioner) sendRegisterRequest(delay time.Duration) {
	if c.requestScheduled {
		return
	}
	c.requestScheduled = true
	if delay > 0 {
		time.Sleep(delay)
	}

	body, err := json.Marshal(RegistrationRequest{
		RegistrationID: c.opts.RegistrationID,
	})
	if err != nil {
		return
	}
	if c.mqttClient != nil {
		// rid := uuid.Must(uuid.NewV4())
		topic := fmt.Sprintf("$dps/registrations/PUT/iotdps-register/?$rid=%s", c.opts.RegistrationID)
		log.WithFields(log.Fields{
			"topic":   topic,
			"payload": string(body),
		}).Info("sending register message")
		c.mqttClient.Publish(topic, 1, false, body)
	} else {
		url := fmt.Sprintf("https://%s/%s/registrations/%s/register", c.opts.Endpoint, c.opts.Scope, c.opts.RegistrationID)
		log.WithFields(log.Fields{
			"url": url,
		}).Info("sending register request")
		c.sendHTTPRequest(http.MethodPut, url, body)
	}
	c.requestScheduled = false
}

func (c *AzureIoTHubProvisioner) sendOperationStatusRequest(delay time.Duration, msg message) {
	if delay > 0 {
		time.Sleep(delay)
	}
	if c.mqttClient != nil {
		// rid := uuid.Must(uuid.NewV4())
		topic := fmt.Sprintf("$dps/registrations/GET/iotdps-get-operationstatus/?$rid=%s&operationId=%s", c.opts.RegistrationID, msg.registrationResponse.OperationID)
		log.WithFields(log.Fields{
			"topic": topic,
		}).Info("sending operation status message")
		c.mqttClient.Publish(topic, 1, false, " ")
	} else {
		url := fmt.Sprintf("https://%s/%s/registrations/%s/operations/%s", c.opts.Endpoint, c.opts.Scope, c.opts.RegistrationID, msg.registrationResponse.OperationID)
		log.WithFields(log.Fields{
			"url": url,
		}).Info("sending operation status request")
		c.sendHTTPRequest(http.MethodGet, url, nil)
	}
}

func (c *AzureIoTHubProvisioner) sendHTTPRequest(method, url string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("api-version", "2019-03-31")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	receivedMsg := &message{
		statusCode: resp.StatusCode,
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.WithField("body", string(respBody)).Info("response")
	json.Unmarshal(respBody, &receivedMsg.registrationResponse)
	if err != nil {
		return err
	}
	// bodyDec := json.NewDecoder(resp.Body)
	// err = bodyDec.Decode(&receivedMsg.registrationResponse)
	// if err != nil && err != io.EOF {
	// 	return err
	// }
	c.messageChan <- *receivedMsg
	return nil
}

func (c *AzureIoTHubProvisioner) messageHandler(client MQTT.Client, msg MQTT.Message) {
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
	if err := json.Unmarshal(msg.Payload(), &receivedMsg.registrationResponse); err != nil {
		statusCode = processingError
	}
	receivedMsg.statusCode = statusCode
	c.messageChan <- receivedMsg
}

func (c *AzureIoTHubProvisioner) writeConfigFile(registrationState RegistrationState) {
	log.WithField("output file", c.opts.OutputFile).Info("writing config file")
	// b, err := ioutil.ReadFile(c.opts.OutputFile)
	// if err == nil {
	// 	viper.SetConfigType("toml")
	// 	viper.MergeConfig(bytes.NewBuffer(b))
	// }
	viper.Set("integration.mqtt.auth.azure_iot_hub.hostname", registrationState.AssignedHub)
	viper.Set("integration.mqtt.auth.azure_iot_hub.device_id", registrationState.DeviceID)
	viper.Set("integration.mqtt.auth.azure_iot_hub.tls_cert", c.opts.Cert)
	viper.Set("integration.mqtt.auth.azure_iot_hub.tls_key", c.opts.Key)
	viper.Set("integration.mqtt.auth.azure_iot_hub.provisioning.Endpoint", c.opts.Endpoint)
	viper.Set("integration.mqtt.auth.azure_iot_hub.provisioning.scope", c.opts.Scope)

	config := viper.AllSettings()
	tree, err := toml.TreeFromMap(config)
	if err != nil {
		log.WithError(err).Error("error creating toml tree")
		return
	}
	file, err := os.OpenFile(c.opts.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	if err != nil {
		log.WithError(err).Error("error opening config file")
		return
	} else {
		if _, err := file.WriteString(tree.String()); err != nil {
			log.WithError(err).Error("error writing config file")
			return
		}
	}
}
