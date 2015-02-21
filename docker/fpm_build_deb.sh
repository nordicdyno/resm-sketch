#!/bin/sh
mkdir -p /root/resm/usr/local/bin
cp /src/bin/resm /root/resm/usr/local/bin/resm
cd /root/resm
fpm -s dir -t deb -v 1.0 -n resm-go .
