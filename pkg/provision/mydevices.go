package provision

import "net/http"

// Client connects to the provisioning server and provisions devices.
type MyDevicesProvioner struct {
	httpClient http.Client
	opts       ProvisionerOptions
}

func NewMyDevicesProvioner(options ProvisionerOptions) *MyDevicesProvioner {
	myd := &MyDevicesProvioner{
		httpClient: http.Client{},
	}

	return myd
}

// ProvisionDevice connects to the server and provisions the device.
func (c *MyDevicesProvioner) ProvisionDevice() {

}
