package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mydevicesiot/dps-client/pkg/provision"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var version string // set by the compiler

func initConfig() provision.Options {
	var opts provision.Options
	var inputFile string
	flag.StringVar(&opts.Endpoint, "e", "global.azure-devices-provisioning.net", "Device provisioning Endpoint URI")
	flag.StringVar(&opts.Scope, "s", "", "Device provisioning scope ID")
	flag.StringVar(&opts.RegistrationID, "r", "", "Registration ID of the device")
	flag.StringVar(&opts.Cert, "c", "", "Full path to the device certificate")
	flag.StringVar(&opts.Key, "k", "", "Full path to the device private key")
	flag.StringVar(&inputFile, "i", "", "Input file containing the opts settings")
	flag.StringVar(&opts.OutputFile, "o", "", "Output file containing the opts settings")
	trace := flag.Bool("t", false, "Output trace info")
	ver := flag.Bool("v", false, "Output version info")
	flag.Parse()

	if *ver == true {
		fmt.Println(version)
		os.Exit(0)
	}

	if *trace == true {
		log.SetLevel(log.TraceLevel)
	}

	if inputFile != "" {
		b, err := ioutil.ReadFile(inputFile)
		if err != nil {
			log.WithError(err).WithField("inputFile", inputFile).Fatal("error loading input file")
		}
		viper.SetConfigType("toml")
		if err := viper.ReadConfig(bytes.NewBuffer(b)); err != nil {
			log.WithError(err).WithField("inputFile", inputFile).Fatal("error loading input file")
		}
	}

	viper.SetDefault("integration.marshaler", "json")
	viper.SetDefault("integration.mqtt.auth.type", "azure_iot_hub")

	viper.SetDefault("integration.mqtt.auth.azure_iot_hub.provisioning.Endpoint", "global.azure-devices-provisioning.net")
	viper.SetDefault("integration.mqtt.auth.azure_iot_hub.provisioning.scope", "0ne00061135")

	viper.SetDefault("commands.commands.reboot.command", "/opt/mydevices/command-ctrl.sh reboot")
	viper.SetDefault("commands.commands.reboot.max_execution_duration", "1s")
	viper.SetDefault("commands.commands.reboot.token", "awVDzEM5S+6dsRJOtF+9lg==")

	viper.SetDefault("commands.commands.remote-ctrl.command", "/opt/mydevices/command-ctrl.sh remote-ctrl")
	viper.SetDefault("commands.commands.remote-ctrl.max_execution_duration", "15s")
	viper.SetDefault("commands.commands.remote-ctrl.token", "RJ5IajLcR5aAbhv/0mvdxw==")

	viper.SetDefault("commands.commands.update.command", "/opt/mydevices/command-ctrl.sh update")
	viper.SetDefault("commands.commands.update.max_execution_duration", "20m")
	viper.SetDefault("commands.commands.update.token", "orrAp+5+RwSW96qWFB4tog==")

	if opts.Endpoint == "" {
		opts.Endpoint = viper.GetString("integration.mqtt.auth.azure_iot_hub.provisioning.Endpoint")
	}
	if opts.Scope == "" {
		opts.Scope = viper.GetString("integration.mqtt.auth.azure_iot_hub.provisioning.scope")
	}
	if opts.RegistrationID == "" {
		opts.RegistrationID = viper.GetString("integration.mqtt.auth.azure_iot_hub.device_id")
	}
	if opts.Cert == "" {
		opts.Cert = viper.GetString("integration.mqtt.auth.azure_iot_hub.tls_cert")
	}
	if opts.Key == "" {
		opts.Key = viper.GetString("integration.mqtt.auth.azure_iot_hub.tls_key")
	}
	log.WithFields(log.Fields{
		"endpoint":       opts.Endpoint,
		"scope":          opts.Scope,
		"registrationID": opts.RegistrationID,
		"cert":           opts.Cert,
		"key":            opts.Key,
		"inputFile":      inputFile,
		"outputFile":     opts.OutputFile,
	}).Trace("options")
	if opts.Endpoint == "" || opts.Scope == "" || opts.RegistrationID == "" || opts.Cert == "" || opts.Key == "" || opts.OutputFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	return opts
}

func main() {
	opts := initConfig()
	client := provision.NewClient(opts)
	client.ProvisionDevice()
}
