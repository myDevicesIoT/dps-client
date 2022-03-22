#!/bin/bash

PACKAGE_NAME="gateway-bridge-dps"
PACKAGE_VERSION=$1
REV="r0"

#Get absolute path with readlink or opkg-build will fail
ROOT_DIR=$(readlink -f "../../")
BUILD_DIR="${ROOT_DIR}/build"
PACKAGE_FILE="${BUILD_DIR}/mipsle/dps-client"
OUTPUT_DIR="${BUILD_DIR}/package/gemtek"
PACKAGE_DIR="${OUTPUT_DIR}/build"
USR_BIN_DIR="${PACKAGE_DIR}/usr/bin/"
TMP_DIR="${PACKAGE_DIR}/tmp/"
APP_LORA_PKG_DIR="${PACKAGE_DIR}/app/lora_pkg/"
# Cleanup
rm -rf $PACKAGE_DIR

# CONTROL
mkdir -p $PACKAGE_DIR/CONTROL
cat > $PACKAGE_DIR/CONTROL/control << EOF
Package: $PACKAGE_NAME
Version: $PACKAGE_VERSION-$REV
Architecture: ramips_24kec
Maintainer: myDevices, Inc. <support@mydevices.com>
Priority: optional
Section: network
Source: N/A
Description: Azure device provisioning client
EOF

cat > $PACKAGE_DIR/CONTROL/postinst << EOF
cp /tmp/chirpstack-gateway-bridge /usr/bin/
/etc/init.d/azure-iot.service stop >> /tmp/gateway-bridge-dps-update.log
/etc/init.d/azure-iot.service start >> /tmp/gateway-bridge-dps-update.log
EOF
chmod 755 $PACKAGE_DIR/CONTROL/postinst

# cat > $PACKAGE_DIR/CONTROL/prerm << EOF
# /etc/init.d/azure-iot.service stop
# EOF
# chmod 755 $PACKAGE_DIR/CONTROL/prerm

cat > $PACKAGE_DIR/CONTROL/conffiles << EOF
EOF

# Files
mkdir -p $USR_BIN_DIR
mkdir -p $TMP_DIR
mkdir -p $APP_LORA_PKG_DIR

cp files/chirpstack-gateway-bridge $TMP_DIR
cp files/lora_wdg_pkt_fwd.sh $APP_LORA_PKG_DIR
cp $PACKAGE_FILE $USR_BIN_DIR

# Package
opkg-build -o root -g root $PACKAGE_DIR $OUTPUT_DIR

# Cleanup
rm -rf $PACKAGE_DIR
