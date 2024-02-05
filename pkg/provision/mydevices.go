package provision

import (
	"net/http"
	"os"

	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Client connects to the provisioning server and provisions devices.
type MyDevicesProvioner struct {
	httpClient http.Client
	opts       ProvisionerOptions
}

func NewMyDevicesProvioner(options ProvisionerOptions) *MyDevicesProvioner {
	myd := &MyDevicesProvioner{
		opts: options,
	}

	return myd
}

// ProvisionDevice connects to the server and provisions the device.
func (c *MyDevicesProvioner) ProvisionDevice() {
	// lets write the file to toml using generic mqtt config
	c.writeConfigFile(RegistrationState{
		AssignedHub:    c.opts.Endpoint,
		RegistrationID: c.opts.RegistrationID,
	})
}

func (c *MyDevicesProvioner) writeConfigFile(registrationState RegistrationState) {
	log.WithField("output file", c.opts.OutputFile).Info("writing config file")
	viper.Set("integration.marshaler", "protobuf")
	viper.Set("integration.mqtt.auth.type", "generic")
	viper.Set("integration.mqtt.auth.generic.hostname", registrationState.AssignedHub)
	viper.Set("integration.mqtt.auth.generic.client_id", registrationState.RegistrationID)
	viper.Set("integration.mqtt.auth.generic.tls_cert", c.opts.Cert)
	viper.Set("integration.mqtt.auth.generic.tls_key", c.opts.Key)

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
