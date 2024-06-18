#!/bin/sh

case "$1" in
    "chirpstack_version")
        /usr/bin/my_devices/chirpstack-gateway-bridge_mydevices version
        ;;
    "dps_client_version")
        /usr/bin/my_devices/dps-client -v
        ;;
    "eth_ip")
        ifconfig eth0 | awk '/inet addr/{print $2}' | cut -d: -f2
        ;;
    "wwan_ip")
        ubus call yruo_status get '{"base":"yruo_status_network"}' | jsonpath -e '$.get[0].value.ip'
        ;;        
    "imei")
        ubus call yruo_status get '{"base":"yruo_celluar"}' |  jsonpath -e '$.get[0].value.imei'
        ;;        
    "mac")
        ifconfig eth0 | awk '/HWaddr/{print $5}'
        ;;
    "eui")
        jsonpath -i /etc/quagga/lora/local_conf.json -e '$.gateway_conf.gateway_ID'
        ;;
    "model")
        ubus call yruo_status get '{"base":"summary"}' | jsonpath -e '$.get[0].value.model'
        ;;
    "serial")
        urtool -g | awk '/^sn/{print $3}'
        ;;
    "firmware_version")
		cat /etc/issue
        ;;
    "marshaler")
		grep "marshaler = " /usr/bin/my_devices/output.toml | awk -F "\"" '{print $2; exit}'
        ;;

    *)
        exit 1
    ;;
esac
