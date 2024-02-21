#!/bin/sh

CONFIG_FILE="/etc/chirpstack-gateway-bridge/chirpstack-gateway-bridge.toml"

case "$1" in
    "cert_expiration")
        CERT_PATH=$(grep tls_cert $CONFIG_FILE | cut -d \" -f2)
        openssl x509 -noout -in $CERT_PATH -enddate  | cut -d= -f2
        ;;
    "chirpstack_version")
        /user/chirpstack-gateway-bridge/chirpstack-gateway-bridge version
        ;;
    "dps_client_version")
        /user/mydevices/dps-client -v
        ;;
    "eth_ip")
        ip addr show eth0 | grep 'inet\b' | awk '{print $2}' | cut -d/ -f1
        ;;
    "wwan_ip")
        # snmpget -v2c -Oqv -m all -c set localhost wmPublicIPAddr.0 2> /dev/null
        ip addr show wwan0 | grep 'inet\b' | awk '{print $2}' | cut -d/ -f1
        ;;        
    "imei")
        snmpget -v2c -Oqv -m all -c set localhost wmImeiNumber.0 2> /dev/null
        ;;
    "imsi")
        snmpget -v2c -Oqv -m all -c set localhost wmImsi.0 2> /dev/null
        ;;        
    "mac")
        cat /sys/class/net/eth0/address
        ;;
    "eui")
        python3 -c 'from keroslib import utils; print(utils.getEUI64())' | tr '[:upper:]' '[:lower:]'
        ;;
    "manufacturer")
        echo "Kerlink"
        ;;
    "model")
        python3 -c 'from keroslib import utils; print(utils.getPlatform().replace("\n", " ") + utils.getRegion())'
        ;;
    "serial")
        python3 -c 'from keroslib import utils; print(utils.getEUI64())' | tr '[:upper:]' '[:lower:]'
        ;;
    "apn")
        gsmdiag.py | awk -F: '/Type:internet/{apn++} apn&&/AccessPointName/{print $2;exit}'
        ;;    
    "firmware_version")
        cat /etc/version
        ;;
    "power_source")
        snmpget -v2c -Oqv -m all -c set localhost pwrSource.0 2> /dev/null
        ;;
    "voltage")
        snmpget -v2c -Oqv -m all -c set localhost pwrInputVoltage.0 2> /dev/null
        ;;
    "battery_critical")
        echo "0"
        ;;
    "charging")
        echo "0"
        ;;
    "charge_fault")
        echo "0"
        ;;                
    "charge_complete")
        echo "0"
        ;;
    "rssi")
        snmpget -v2c -Oqv -m all -c set localhost wmSignalStrength.0 2> /dev/null
        ;;        
    "marshaler")
        grep marshaler $CONFIG_FILE | grep -v $(basename "$0") | cut -d \" -f2
        ;;       
    *)
        exit 1
    ;;
esac
