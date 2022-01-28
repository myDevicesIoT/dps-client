# Device Provisioning Service Client

The Device Provisioning Service Client provisions devices with an Azure IoT Hub using the Azure Device Provisioning Service and creates/updates the [chirpstack-gateway-bridge](https://github.com/brocaar/chirpstack-gateway-bridge) config file with the Azure IoT Hub connection information.

# Build

To build for the current OS:

```
make
```

To cross compile for a specific architecture:

```
make <TARGET_ARCH>
```
Where `TARGET_ARCH` is one of the following: `armv5`, `armv7`, `mips`, `mipsle`.

# Package

To create an opkg package for a specific gateway you can use the provided makefile targets. This requires [opkg-utils](https://git.yoctoproject.org/opkg-utils/) to be installed locally. This can be done by cloning the opkg-utils repo and running `make install`.

```
make <MANUFACTURER>
```
Where `MANUFACTURER` is one of the following: `multitech`, `tektelic`.

# Usage

The dps-client can provision devices using options from a specified input config file. It is recommended to use a default.toml file like the ones under the `packaging/<manufacturer>/files` folders and specify the device ID and cert file paths on the command line.

```
dps-client -i /path/to/default.toml -r <DEVICE_ID> -c /path/to/device/cert.pem -k /path/to/device/key.pem -o /path/to/output_config.toml
```
