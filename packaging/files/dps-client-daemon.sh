#!/bin/bash

CERT_DIR="/var/config/app/mydevices"
DEVEUI=$(mts-io-sysfs show lora/eui | tr -d : | tr '[:upper:]' '[:lower:]')
DEFAULT_CONFIG_FILE="/etc/opt/mydevices/default.toml"
CONFIG_FILE="/var/config/chirpstack-gateway-bridge/chirpstack-gateway-bridge.toml"
TEMP_CONFIG_FILE="$CONFIG_FILE.tmp"
DPS_CLIENT_ARGS="-i $DEFAULT_CONFIG_FILE -r $DEVEUI -c $CERT_DIR/$DEVEUI.cert.pem -k $CERT_DIR/$DEVEUI.key.pem -o $TEMP_CONFIG_FILE"

while :
do
    /opt/mydevices/dps-client $DPS_CLIENT_ARGS
    if cmp -s "$TEMP_CONFIG_FILE" "$CONFIG_FILE" ; then
        rm "$TEMP_CONFIG_FILE"
    else
        mv "$TEMP_CONFIG_FILE" "$CONFIG_FILE"
        /etc/init.d/chirpstack-gateway-bridge restart
    fi
	sleep 7d
done
