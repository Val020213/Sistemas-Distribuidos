#!/bin/sh
iptables -t nat -A POSTROUTING -s 10.0.10.0/24 -o eth0 -j MASQUERADE
iptables -t nat -A POSTROUTING -s 10.0.11.0/24 -o eth0 -j MASQUERADE
go run ./cmd/api/ &
wait $!