#!/bin/sh

# Configurar la ruta por defecto
ip route del default
ip route add default via 10.0.10.254

sysctl -w net.ipv4.ip_forward=1
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE

# Ejecutar el servidor en modo desarrollo
go run ./cmd/api/ 