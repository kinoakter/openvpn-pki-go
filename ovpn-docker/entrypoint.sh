#!/usr/bin/env bash

tun_dev="udp0"
args="--config /etc/openvpn/${tun_dev}.conf --script-security 2"

exec openvpn --config /etc/openvpn/${tun_dev}.conf