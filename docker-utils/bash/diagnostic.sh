#!/bin/bash
CONTAINER_NAME=$1

# Verificar estado del contenedor
docker inspect $CONTAINER_NAME --format '
Estado: {{.State.Status}}
Health: {{.State.Health.Status}}
IP: {{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}
Puertos: {{range $p, $conf := .NetworkSettings.Ports}}{{$p}} {{end}}'

# Probar conectividad interna
echo -e "\nProbando conexión desde el contenedor:"
docker exec $CONTAINER_NAME curl -I http://localhost:3000

# Verificar iptables del router
echo -e "\nReglas de NAT en el router:"
docker exec scrapper-router iptables -t nat -L -n -v

# Probar acceso desde el host
HOST_PORT=$(docker inspect $CONTAINER_NAME --format '{{(index (index .NetworkSettings.Ports "3000/tcp") 0).HostPort}}')
echo -e "\nPrueba desde el host:"
curl -v http://localhost:$HOST_PORT