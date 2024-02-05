package provision

import "context"

// Options contains the device info used for provisioning.
type ProvisionerOptions struct {
	Endpoint       string
	Scope          string
	RegistrationID string
	Cert           string
	Key            string
	OutputFile     string
	Protocol       string
}

// DeviceProvisioner is an interface that describes a device provisioner.
type DeviceProvisioner interface {
	ProvisionDevice(ctx context.Context, deviceID string) error
}
