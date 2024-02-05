package hub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HubClient struct {
	httpClient http.Client
	url        string
}

var mydevicesHub = "global.hub.mydevices.com"

type HubResponse struct {
	DpsClient string `json:"dps-client"`
	CGB       string `json:"chirpstack-gateway-bridge"`
	Endpoint  string `json:"endpoint"`
	Provider  string `json:"provider"`
}

func NewHubClient(gatewayId string) *HubClient {

	hub := &HubClient{}

	httpClient := http.Client{}

	url := fmt.Sprintf("https://%s/api/gateways/%s", mydevicesHub, gatewayId)

	hub.url = url
	hub.httpClient = httpClient
	return hub
}

func (h *HubClient) PingHome() (HubResponse, error) {
	var hubResponse HubResponse
	req, err := http.NewRequest("GET", h.url, nil)
	if err != nil {
		return hubResponse, err
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return hubResponse, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return hubResponse, err
	}
	// body should return a json with the latest version of
	// dps-client, chirpstack, broker, and endpoint.
	// parse the response body and set the values to the HubResponse struct
	err = json.Unmarshal(body, &hubResponse)
	if err != nil {
		return hubResponse, err
	}

	return hubResponse, nil
}
