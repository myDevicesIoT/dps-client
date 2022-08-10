#!/bin/sh

COMMISSION_DB="/tmp/commissioning.db"

case "$1" in
    "cert_expiration")
        CONFIG_FILE="/etc/chirpstack-gateway-bridge/chirpstack-gateway-bridge.toml"
        CERT_PATH=$(grep tls_cert $CONFIG_FILE | cut -d \" -f2)
        openssl x509 -noout -in $CERT_PATH -enddate  | cut -d= -f2
        ;;
    "chirpstack_version")
        /opt/chirpstack-gateway-bridge/chirpstack-gateway-bridge version
        ;;
    "dps_client_version")
        /opt/mydevices/dps-client -v
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
        sqlite3 "${COMMISSION_DB}" "SELECT Customer_Gateway_ID FROM Configuration" | tr '[:upper:]' '[:lower:]'
        ;;
    "manufacturer")
        # snmpget -v2c -Oqv -m all -c set localhost swVerComponentVer.1 2> /dev/null
        system_version | grep "Distributor ID" | awk '{ sub(/^[^:]*:[[:blank:]]*/, ""); print }'
        ;;
    "model")
        # sqlite3 "${COMMISSION_DB}" "SELECT Module_Name FROM Configuration"
        # snmpget -v2c -Oqv -m all -c set localhost swVerComponentVer.4 2> /dev/null
        system_version | grep "Product" | awk '{ sub(/^[^:]*:[[:blank:]]*/, ""); print }'
        ;;
    "serial")
        sqlite3 "${COMMISSION_DB}" "SELECT Module_Serial_Number FROM Configuration"
        ;;
    "apn")
        snmpget -v2c -Oqv -m all -c set localhost wmApnName.0 2> /dev/null
        ;;    
    "firmware_version")
        # snmpget -v2c -Oqv -m all -c set localhost swVerComponentVer.3 2> /dev/null
        system_version | grep "Release" | awk '{ sub(/^[^:]*:[[:blank:]]*/, ""); print }'
        ;;
    "power_source")
        snmpget -v2c -Oqv -m all -c set localhost pwrSource.0 2> /dev/null
        ;;
    "voltage")
        snmpget -v2c -Oqv -m all -c set localhost pwrInputVoltage.0 2> /dev/null
        ;;
    "battery_critical")
        snmpget -v2c -Oqv -m all -c set localhost pwrBatteryCritical.0 2> /dev/null
        ;;
    "charging")
        snmpget -v2c -Oqv -m all -c set localhost pwrBatteryCharging.0 2> /dev/null
        ;;
    "charge_fault")
        snmpget -v2c -Oqv -m all -c set localhost pwrBatteryChargeFault.0 2> /dev/null
        ;;                
    "charge_complete")
        snmpget -v2c -Oqv -m all -c set localhost pwrBatteryChargeComplete.0 2> /dev/null
        ;;
    *)
        exit 1
    ;;
esac
