#!/bin/bash


# ANSI color codes
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Valores por defecto
DEFAULT_NETWORK_NAME="scrapper-server-network"
DEFAULT_SUBNET="10.0.0.0/24"
DEFAULT_MONGO_IMAGE="mongo:6.0"
DEFAULT_SERVER_IMAGE="scrapper_server"

# Valores por defecto para MongoDB
MONGO_VOLUME="mongo_data"
MONGO_PORT=27018
MONGO_USERNAME="admin"
MONGO_PASSWORD="password"

# Leer parámetros o usar valores por defecto
while getopts ":m:i:s:j:p:" opt; do
  case $opt in
    m) MONGO_CONTAINER_NAME="$OPTARG"
    ;;
    i) MONGO_IP="$OPTARG"
    ;;
    s) SERVER_CONTAINER_NAME="$OPTARG"
    ;;
    j) SERVER_IP="$OPTARG"
    ;;
    p) SERVER_PORT="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done
# Validar que los parámetros requeridos estén presentes
if [ -z "$MONGO_CONTAINER_NAME" ] || [ -z "$MONGO_IP" ] || [ -z "$SERVER_CONTAINER_NAME" ] || [ -z "$SERVER_IP" ] || [ -z "$SERVER_PORT" ]; then
    echo "Faltan parámetros requeridos. Uso: $0 -m <mongo_container_name> -i <mongo_ip> -s <server_container_name> -j <server_ip> -p <server_port>"
    exit 1
fi
NETWORK_NAME=$DEFAULT_NETWORK_NAME
MONGO_CONTAINER_NAME=${MONGO_CONTAINER_NAME:-"mongo"}
MONGO_IP=${MONGO_IP:-"10.0.10.3"}
SERVER_CONTAINER_NAME=${SERVER_CONTAINER_NAME:-"scrapper_server"}
SERVER_IP=${SERVER_IP:-"10.0.10.2"}
SERVER_PORT=${SERVER_PORT:-8080}


# Crear volumen para MongoDB si no existe
docker volume create $MONGO_VOLUME
docker run -d \
    --name $MONGO_CONTAINER_NAME \
    --rm \
    --network $NETWORK_NAME \
    --ip $MONGO_IP \
    -p $MONGO_PORT:$MONGO_PORT \
    -e MONGO_INITDB_ROOT_USERNAME=$MONGO_USERNAME \
    -e MONGO_INITDB_ROOT_PASSWORD=$MONGO_PASSWORD \
    -v $MONGO_VOLUME:/data/db \
    $DEFAULT_MONGO_IMAGE

echo -e "${GREEN}MongoDB iniciado en el contenedor '$MONGO_CONTAINER_NAME' con IP $MONGO_IP${NC}"

# Construir la imagen del servidor desde el Dockerfile
if ! docker image inspect "$DEFAULT_SERVER_IMAGE" > /dev/null 2>&1; then
    echo -e "${GREEN}Imagen $DEFAULT_SERVER_IMAGE no encontrada localmente. Construyendo la imagen...${NC}"
    docker build -t "$DEFAULT_SERVER_IMAGE" ./scrapper_server || { echo -e "${GREEN}Error al construir la imagen $DEFAULT_SERVER_IMAGE.${NC}"; exit 1; }
fi

# Iniciar el servidor en la red personalizada con una IP fija
docker run -d \
    --name $SERVER_CONTAINER_NAME \
    --rm \
    --network $NETWORK_NAME \
    --ip $SERVER_IP \
    -p $SERVER_PORT:8080 \
    -e PORT=$SERVER_PORT \
    -e BLUEPRINT_DB_HOST=$MONGO_IP \
    -e BLUEPRINT_DB_PORT=$MONGO_PORT \
    -e BLUEPRINT_DB_USER=$MONGO_USERNAME \
    -e BLUEPRINT_DB_PASSWORD=$MONGO_PASSWORD \
    -v $(pwd):/app \
    -v /app/go/pkg \
    $DEFAULT_SERVER_IMAGE

echo -e "${GREEN}Servidor iniciado en el contenedor '$SERVER_CONTAINER_NAME', escuchando en el puerto $SERVER_PORT con IP $SERVER_IP${NC}"
