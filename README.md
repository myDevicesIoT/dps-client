# Device Provisioning Service Client

The Device Provisioning Service Client provisions devices with an Azure IoT Hub using the Azure Device Provisioning Service and creates/updates the [chirpstack-gateway-bridge](https://github.com/brocaar/chirpstack-gateway-bridge) config file with the Azure IoT Hub connection information.

# Build with Docker (Recommended)

To build and package the `dps-client` binary using Docker, run the following command:
```shell
docker run -v $PWD:/opt/dps-client -it mydevices/dps-client:dev-build /bin/bash
```

This will mount the current directory to the `/opt/dps-client` directory in the container. The `dps-client` source code is located in the `/opt/dps-client` directory in the container.

To build the `dps-client` binary, run the following command:

```shell
make multitech
```


# Build Locally

To build for the current OS:

```
make
```

To cross compile for a specific architecture:

```
make <TARGET_ARCH>
```
Where `TARGET_ARCH` is one of the following: `armv5`, `armv7`, `mips`, `mipsle`.

This command will build the `dps-client` at `build/<TARGET_ARCH>`. The binary can be run with the parameters below and will output a config file if provisioning is successful. The config file is used with the `chirpstack-gateway-bridge` binary.

# Package

To create an opkg package for a specific gateway you can use the provided makefile targets. This requires [opkg-utils](https://git.yoctoproject.org/opkg-utils/) to be installed locally. This can be done by cloning the opkg-utils repo and running `make install` from the opkg-utils root directory.

```
make <MANUFACTURER>
```
Where `MANUFACTURER` is one of the following: `kerlink`, `multitech`, `tektelic` or `gemtek`. The package will be created at `build/package/<MANUFACTURER>`.

# Docker
A *Dockerfile* is provided for ease of compiling on other architectures

1. build image

```shell
docker build -t mydevices/dps-client:latest .
```
2. run headless container

Example of how to run the binary
```shell
docker run -w ./ -i -t mydevices/dps-client:latest /opt/dps-client/dps-client --help 
```

```
docker run -v $(pwd):'/home/' -i -t mydevices/dps-client:latest /opt/dps-client/dps-client -e global.azure-devices-provisioning.net -s 0ne0006B4D6 -r GATEWAY_ID -c /home/GATEWAY_ID.cert.pem -k /home/GATEWAY_ID.key.pem  -o /home/chirpstack-gateway-bridge.toml
```
The above command assumes that the current host has the certificates in a local directory. (pwd) 



# Usage

The `dps-client` can provision devices using options from a specified input config file. It is recommended to use a `default.toml` file like the ones under the `packaging/<manufacturer>/files` folders and specify the device ID and cert file paths on the command line. If using an input file it must contain the `[integration]` section like the ones under the `packaging/<manufacturer>/files` folders.

```
[integration]
  marshaler = "json"

  [integration.mqtt]

    [integration.mqtt.auth]
      type = "azure_iot_hub"

        [integration.mqtt.auth.azure_iot_hub.provisioning]
          endpoint = "global.azure-devices-provisioning.net"
          scope = "0ne0006B4D6"
```

If support for gateway commands (reboot, remote-ctrl and update) are needed a `[commands]` section is also required and a `command-ctrl.sh` script will need to be created for the target device. Examples of the `command-ctrl.sh` script are also under the `packaging/<manufacturer>/files` folders.

```
[commands]

  [commands.commands]

    [commands.commands.reboot]
      command = "/opt/mydevices/command-ctrl.sh reboot"
      max_execution_duration = "1s"

    [commands.commands.remote-ctrl]
      command = "/opt/mydevices/command-ctrl.sh remote-ctrl"
      max_execution_duration = "15s"

    [commands.commands.update]
      command = "/opt/mydevices/command-ctrl.sh update"
      max_execution_duration = "20m"
```

Command with an input toml file specified: 

```
dps-client -i /path/to/default.toml -r <DEVICE_ID> -c /path/to/device/cert.pem -k /path/to/device/key.pem -o /path/to/output_config.toml
```

Alternatively the `dps-client` can be run without an input toml file by specifying the endpoint and scope id on the command line. In this case a standard default config file will be generated. The available command line options are specified below.

```
Usage of ./dps-client:
  -c string
        Full path to the device certificate
  -e string
        Device provisioning Endpoint URI (default "global.azure-devices-provisioning.net")
  -i string
        Input file containing the opts settings
  -k string
        Full path to the device private key
  -o string
        Output file containing the opts settings
  -p string
        Protocol to use when provisioning device, mqtt or https (default "mqtt")
  -r string
        Registration ID of the device
  -s string
        Device provisioning scope ID
  -t    Output trace info
  -v    Output version info
```

# Packages

The current versions of packages for the supported gateways are available here:

- [Multitech](https://hwdartifacts.blob.core.windows.net/hwdassets/dps-client_1.3.8-r0_arm926ejste.ipk)
- [Tektelic](https://hwdartifacts.blob.core.windows.net/hwdassets/dps-client_1.3.8-r0_kona.ipk)
- [Gemtek](https://hwdartifacts.blob.core.windows.net/hwdassets/gateway-bridge-dps_1.3.9-r0_ramips_24kec.ipk)
