#!/bin/bash

PACKAGE_NAME="dps-client"
PACKAGE_VERSION=$1
REV="r0"

#Get absolute path with readlink or opkg-build will fail
ROOT_DIR=$(readlink -f "../../")
BUILD_DIR="${ROOT_DIR}/build"
PACKAGE_FILE="${BUILD_DIR}/armv5/dps-client"
OUTPUT_DIR="${BUILD_DIR}/package/tektelic"
PACKAGE_DIR="${OUTPUT_DIR}/build"
OPT_DIR="${PACKAGE_DIR}/opt/mydevices"
ETC_OPT_DIR="${PACKAGE_DIR}/etc/opt/mydevices"
INIT_DIR="${PACKAGE_DIR}/etc/init.d"
MONIT_DIR="${PACKAGE_DIR}/etc/monit.d"

# Cleanup
rm -rf $PACKAGE_DIR

# CONTROL
mkdir -p $PACKAGE_DIR/CONTROL
cat > $PACKAGE_DIR/CONTROL/control << EOF
Package: $PACKAGE_NAME
Version: $PACKAGE_VERSION-$REV
Architecture: kona
Maintainer: myDevices, Inc. <support@mydevices.com>
Priority: optional
Section: network
Source: N/A
Description: Azure device provisioning client
EOF

cat > $PACKAGE_DIR/CONTROL/postinst << EOF
/usr/bin/monit reload
EOF
chmod 755 $PACKAGE_DIR/CONTROL/postinst

cat > $PACKAGE_DIR/CONTROL/prerm << EOF
/etc/init.d/dps-client stop
EOF
chmod 755 $PACKAGE_DIR/CONTROL/prerm

cat > $PACKAGE_DIR/CONTROL/conffiles << EOF
/etc/opt/mydevices/default.toml
EOF

# Files
mkdir -p $OPT_DIR
mkdir -p $ETC_OPT_DIR
mkdir -p $INIT_DIR
mkdir -p $MONIT_DIR

cp files/$PACKAGE_NAME.init $INIT_DIR/$PACKAGE_NAME
cp files/$PACKAGE_NAME.monit $MONIT_DIR/$PACKAGE_NAME
cp files/command-ctrl.sh $OPT_DIR
chmod 755 $OPT_DIR/command-ctrl.sh
cp files/meta-data.sh $OPT_DIR
chmod 755 $OPT_DIR/meta-data.sh
cp files/dps-client-daemon.sh $OPT_DIR
cp $PACKAGE_FILE $OPT_DIR
cp files/default.toml $ETC_OPT_DIR
chmod 755 $ETC_OPT_DIR/default.toml

# Package
opkg-build -o root -g root $PACKAGE_DIR $OUTPUT_DIR

# Cleanup
rm -rf $PACKAGE_DIR
