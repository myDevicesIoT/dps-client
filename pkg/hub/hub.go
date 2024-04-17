package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HubClient struct {
	httpClient http.Client
	url        string
}

var mydevicesHub = "global-provisioning.mydevices.com"

type HubResponse struct {
	GatewayID string `json:"id"`
	Endpoint  string `json:"endpoint"`
	Provider  string `json:"provider"`
	Command   string `json:"command"`
}

func NewHubClient(gatewayId string, dpsVersion string) *HubClient {

	hub := &HubClient{}

	httpClient := http.Client{
		Timeout: 5 * time.Second * 5,
	}

	url := fmt.Sprintf("https://%s/api/gateways/%s/connect?dps_version=%s", mydevicesHub, gatewayId, dpsVersion)

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
		return hubResponse, errors.New("error reading mydevices reading")
	}
	// body should return a json with the latest version of
	// dps-client, chirpstack, broker, and endpoint.
	// parse the response body and set the values to the HubResponse struct
	err = json.Unmarshal(body, &hubResponse)
	if err != nil {
		return hubResponse, errors.New("error parsing mydevices response")
	}

	return hubResponse, nil
}
