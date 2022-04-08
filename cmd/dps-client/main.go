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

var version string                                       // set by the compiler
var commandScriptPath = "/opt/mydevices/command-ctrl.sh" // can be overridden by the compiler

func initConfig() provision.Options {
	var opts provision.Options
	var inputFile string
	flag.StringVar(&opts.Endpoint, "e", "", "Device provisioning Endpoint URI")
	flag.StringVar(&opts.Scope, "s", "", "Device provisioning scope ID")
	flag.StringVar(&opts.RegistrationID, "r", "", "Registration ID of the device")
	flag.StringVar(&opts.Cert, "c", "", "Full path to the device certificate")
	flag.StringVar(&opts.Key, "k", "", "Full path to the device private key")
	flag.StringVar(&inputFile, "i", "", "Input file containing the opts settings")
	flag.StringVar(&opts.OutputFile, "o", "", "Output file containing the opts settings")
	flag.StringVar(&opts.Protocol, "p", "mqtt", "Protocol to use when provisioning device, mqtt or https")
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

	rebootCommand := fmt.Sprintf("%s reboot", commandScriptPath)
	viper.SetDefault("commands.commands.reboot.command", rebootCommand)
	viper.SetDefault("commands.commands.reboot.max_execution_duration", "1s")

	remoteCommand := fmt.Sprintf("%s remote-ctrl", commandScriptPath)
	viper.SetDefault("commands.commands.remote-ctrl.command", remoteCommand)
	viper.SetDefault("commands.commands.remote-ctrl.max_execution_duration", "15s")

	updateCommand := fmt.Sprintf("%s update", commandScriptPath)
	viper.SetDefault("commands.commands.update.command", updateCommand)
	viper.SetDefault("commands.commands.update.max_execution_duration", "20m")

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
	log.SetOutput(os.Stdout)
	opts := initConfig()
	client := provision.NewClient(opts)
	client.ProvisionDevice()
}
