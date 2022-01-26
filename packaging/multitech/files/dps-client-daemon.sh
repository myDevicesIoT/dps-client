#!/bin/bash

CERT_DIR="/var/config/app/mydevices"
OLD_CERT_DIR="/var/config/mydevices"
DEVEUI=$(mts-io-sysfs show lora/eui | tr -d : | tr '[:upper:]' '[:lower:]')
DEFAULT_CONFIG_FILE="/etc/opt/mydevices/default.toml"
CONFIG_FILE="/var/config/chirpstack-gateway-bridge/chirpstack-gateway-bridge.toml"
TEMP_CONFIG_FILE="$CONFIG_FILE.tmp"
DPS_CLIENT_ARGS="-i $DEFAULT_CONFIG_FILE -r $DEVEUI -c $CERT_DIR/$DEVEUI.cert.pem -k $CERT_DIR/$DEVEUI.key.pem -o $TEMP_CONFIG_FILE"

move_cert() {
    #Move cert files to the new location, if they aren't already there.
    CERT_FILE=$1
    if [ ! -f "$CERT_DIR/$CERT_FILE" ]; then
        echo "$CERT_DIR/$CERT_FILE does not exist."
        mkdir -p $CERT_DIR
        mv $OLD_CERT_DIR/$CERT_FILE $CERT_DIR/
    fi
}

move_cert $DEVEUI.cert.pem
move_cert $DEVEUI.key.pem

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
