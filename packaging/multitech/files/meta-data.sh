#!/bin/sh


case "$1" in
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
        curl -s localhost/api/stats/eth0 | python -c "import sys, json; print json.load(sys.stdin)['result']['ip']"
        ;;
    "ppp_ip")
        curl -s localhost/api/stats/ppp | python -c "import sys, json; print json.load(sys.stdin)['result']['localIp']"
        ;;
    "apn")
        curl -s localhost/api/ppp/modem | python -c "import sys, json; print json.load(sys.stdin)['result']['apnString']"
        ;;
    "today_tx")
        curl -s localhost/api/stats/pppTotal | python -c "import sys, json; print json.load(sys.stdin)['result']['todayTx']"
        ;;
    "today_rx")
        curl -s localhost/api/stats/pppTotal | python -c "import sys, json; print json.load(sys.stdin)['result']['todayRx']"
        ;;
    "firmware_version")
        curl -s localhost/api/system | python -c "import sys, json; print json.load(sys.stdin)['result']['firmware']"
        ;;
    *)
        exit 1
    ;;
esac
