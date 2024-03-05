#!/bin/sh

CONFIG_FILE="/var/config/chirpstack-gateway-bridge/chirpstack-gateway-bridge.toml"

function get_json_value () {
    endpoint=$1
    key=$2
    if [ -x "$(command -v python3)" ]; then
        echo $(curl -s $endpoint | python3 -c "import sys, json; print(json.load(sys.stdin)['result']['$key'])")
    else
        echo $(curl -s $endpoint | python -c "import sys, json; print json.load(sys.stdin)['result']['$key']")
    fi
}

case "$1" in
    "cert_expiration")
        CERT_PATH=$(grep tls_cert $CONFIG_FILE | cut -d \" -f2)
        openssl x509 -noout -in $CERT_PATH -enddate  | cut -d= -f2
        ;;
    "chirpstack_version")
        /opt/chirpstack-gateway-bridge/chirpstack-gateway-bridge version
        ;;
    "dps_client_version")
        /opt/mydevices/dps-client -v
        ;;
    "eui")
        mts-io-sysfs show lora/eui | tr -d : | tr '[:upper:]' '[:lower:]'
        ;;    
    "eth_ip")
        get_json_value "localhost/api/stats/eth0" "ip"
        ;;
    "ppp_ip")
        get_json_value "localhost/api/stats/ppp" "localIp"
        ;;
    "apn")
        get_json_value "localhost/api/ppp/modem" "apnString"
        ;;
    "imsi")
        get_json_value "localhost/api/system" "imsi"
        ;;
    "today_tx")
        get_json_value "localhost/api/stats/pppTotal" "todayTx"
        ;;
    "today_rx")
        get_json_value "localhost/api/stats/pppTotal" "todayRx"
        ;;
    "firmware_version")
        get_json_value "localhost/api/system" "firmware"
        ;;
    "rssi")
        get_json_value "localhost/api/stats/ppp" "rssiDbm"
        ;;
    "marshaler")
        grep marshaler $CONFIG_FILE | grep -v $(basename "$0") | cut -d \" -f2
        ;;
    *)
        exit 1
    ;;
esac
