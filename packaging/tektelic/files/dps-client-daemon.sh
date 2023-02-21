#!/bin/bash

CERT_DIR="/var/config/mydevices"
DEVEUI=""
COMMISSION_DB="/tmp/commissioning.db"
DEFAULT_CONFIG_FILE="/etc/opt/mydevices/default.toml"
CONFIG_DIR="/etc/chirpstack-gateway-bridge"
CONFIG_FILE="$CONFIG_DIR/chirpstack-gateway-bridge.toml"
TEMP_CONFIG_FILE="$CONFIG_FILE.tmp"
OLD_CONFIG_FILE="/var/config/chirpstack-gateway-bridge/chirpstack-gateway-bridge.toml"

obtain_gwid() {
    for ((i=0; i<60; i++)); do
        if [ -r "${COMMISSION_DB}" ]; then
            DEVEUI=$(sqlite3 "${COMMISSION_DB}" "SELECT Customer_Gateway_ID FROM Configuration" | tr '[:upper:]' '[:lower:]')
            return 0
        fi
        echo "Commissioning database not available yet"
        sleep 1
    done

    return 1
}

obtain_gwid

while :
do
    mkdir -p $CONFIG_DIR
    DPS_CLIENT_ARGS="-i $DEFAULT_CONFIG_FILE -r $DEVEUI -c $CERT_DIR/$DEVEUI.cert.pem -k $CERT_DIR/$DEVEUI.key.pem -o $TEMP_CONFIG_FILE"
    /opt/mydevices/dps-client $DPS_CLIENT_ARGS
    RESTART_CHIRPSTACK=false
    if test -f "$OLD_CONFIG_FILE"; then
        if ! cmp -s "$TEMP_CONFIG_FILE" "$OLD_CONFIG_FILE" ; then
            cp "$TEMP_CONFIG_FILE" "$OLD_CONFIG_FILE"
            RESTART_CHIRPSTACK=true
        fi
    fi
    if cmp -s "$TEMP_CONFIG_FILE" "$CONFIG_FILE" ; then
        rm "$TEMP_CONFIG_FILE"
    else
        mv "$TEMP_CONFIG_FILE" "$CONFIG_FILE"
        RESTART_CHIRPSTACK=true
    fi
    if [ "$RESTART_CHIRPSTACK" = true ]; then
        /etc/init.d/chirpstack-gateway-bridge restart
    fi
    sleep 7d
done
