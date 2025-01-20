#!/bin/sh

# Configurar la ruta por defecto
ip route del default
ip route add default via 10.0.10.254

cd ./cmd/api

# Ejecutar el servidor en modo desarrollo
go run .
