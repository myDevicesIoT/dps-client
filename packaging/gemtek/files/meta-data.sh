#!/bin/sh


case "$1" in
    "cert_expiration")
        CONFIG_FILE="/mnt/data/app/azureiot/chirpstack-gateway-bridge.toml"
        CERT_PATH=$(grep tls_cert $CONFIG_FILE | cut -d \" -f2)
        openssl x509 -noout -in $CERT_PATH -enddate  | cut -d= -f2
        ;;
    "chirpstack_version")
        /mnt/data/myd/myd/chirpstack-gateway-bridge version
        ;;
    "dps_client_version")
        /mnt/data/myd/myd/dps-client -v
        ;;           
    "eth_ip")
        # DEFAULT_ROUTE=$(ip route show default | awk '/default/ {print $5}')
        # ip addr show $DEFAULT_ROUTE | grep 'inet\b' | awk '{print $2}' | cut -d/ -f1
        ip addr show eth0 | grep 'inet\b' | awk '{print $2}' | cut -d/ -f1
        ;;
    "mac")
        # DEFAULT_ROUTE=$(ip route show default | awk '/default/ {print $5}')
        # cat /sys/class/net/$DEFAULT_ROUTE/address
        uci get -c /mnt/data/config mfg.system.base_mac
        ;;
    "eui")
        uci get -c /mnt/data/config mfg.system.lora_ee_cnt_0 | cut -c 1-16
        ;;
    "manufacturer")
        grep MANUFACTUER /mnt/data/config/*.ini | sed 's/^.*"\(.*\)".*/\1/'
        ;;
    "model")
        uci get -c /mnt/data/config mfg.system.model_name
        ;;
    "serial")
        uci get -c /mnt/data/config mfg.system.sn
        ;;
    "hardware_version")
        uci get -c /etc/config profile.system.hw_version
        ;;        
    "firmware_version")
        uci get -c /etc/config profile.system.fw_version
        ;;
    *)
        exit 1
    ;;
esac
