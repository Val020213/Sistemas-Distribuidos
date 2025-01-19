cd ./src || exit 1

export CLIENT_NETWORK="scrapper-client-network"
export SERVER_NETWORK="scrapper-server-network"
export ROUTER_IMAGE="scrapper-router-image"
export ROUTER_CONTAINER="scrapper-router"
export ROUTER_HOSTNAME="scrapper-router"
export ROUTER_ALIAS="scrapper-router"
export ROUTER_SUBNET_SERVER="10.0.10.0/24"
export ROUTER_SUBNET_CLIENT="10.0.11.0/24"
export ROUTER_IP_SERVER="10.0.10.254"
export ROUTER_IP_CLIENT="10.0.11.254"
export ROUTER_VOLUME="./router:/app"

docker images | grep -q "$ROUTER_IMAGE" || {
    echo "Construyendo la imagen $ROUTER_IMAGE..."
    docker-compose build router
}

echo "Iniciando la red y router..."
if docker network ls | grep -q "$SERVER_NETWORK"; then
    docker network inspect "$SERVER_NETWORK" -f '{{range .Containers}}{{.Name}} {{end}}' | xargs -r docker network disconnect "$SERVER_NETWORK"
    docker network rm "$SERVER_NETWORK"
fi
docker network create --subnet=$ROUTER_SUBNET_SERVER "$SERVER_NETWORK"

docker rm -f $ROUTER_CONTAINER 2>/dev/null || true

docker run -d --name $ROUTER_CONTAINER --hostname $ROUTER_HOSTNAME \
    --network $SERVER_NETWORK --ip $ROUTER_IP_SERVER \
    --cap-add=NET_ADMIN --privileged \
    --sysctl net.ipv4.ip_forward=1 \
    -v $ROUTER_VOLUME $ROUTER_IMAGE:latest

docker network connect --ip $ROUTER_IP_CLIENT $CLIENT_NETWORK $ROUTER_CONTAINER
