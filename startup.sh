#!/bin/sh
while ! ping -c1 8.8.8.8 &>/dev/null; do echo "waiting for net"; sleep 1; done
$GOPATH/bin/slacksoc /etc/slacksoc.yaml
