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


# Usage

The dps-client provisions devices by specifying the Azure IoT Hub and device options on the command line:

```
dps-client -r <DEVICE_ID> -s <AZURE_DPS_SCOPE_ID> -c /path/to/device/cert.pem -k /path/to/device/key.pem -o /path/to/config.toml
```

The dps-client can also provision devices using options from a specified input config file:

```
dps-client -i /path/to/config.toml -o /path/to/config.toml
```