# Este script crea y ejecuta un contenedor Docker con una configuración específica.
# 
# Uso:
#   ./add-client.sh -n <nombre_contenedor> -i <ip> [-p <puerto_host:puerto_contenedor>] [-e <env_var=valor>] [-network <red>] [-image <imagen>]
#
# Parámetros:
#   -n, --name       Nombre del contenedor (obligatorio)
#   -i, --ip         Dirección IP del contenedor (obligatorio)
#   -p, --port       Mapeo de puertos en formato <puerto_host:puerto_contenedor> (opcional, por defecto "3000:3000")
#   -e, --env        Variables de entorno en formato <env_var=valor> (opcional, se pueden especificar múltiples)
#   -network, --network Red de Docker a la que se conectará el contenedor (opcional, por defecto $CLIENT_NETWORK)
#   -image, --image  Imagen de Docker a utilizar (opcional, por defecto "scrapper-client-image:latest")
#
# Ejemplo de uso:
#   ./add-client.sh -n mi_contenedor -i 172.18.0.22 -p 8080:80 -e NODE_ENV=production -network mi_red -image mi_imagen:latest
#
# Descripción:
#   Este script inicializa una red Docker, valida los parámetros proporcionados, asigna valores por defecto si es necesario,
#   construye el comando `docker run` con las opciones especificadas y ejecuta el comando para crear y ejecutar el contenedor.
#   Si el contenedor se crea exitosamente, muestra un mensaje de éxito; de lo contrario, muestra un mensaje de error.
#!/bin/bash
#!/bin/bash

DEFAULT_IMAGE="scrapper-client-image:latest"
DEFAULT_VOLUMES=(
  "$(pwd)/client:/app"
  "/app/node_modules"
  "/app/.next"
)
DEFAULT_PORT="3000:3000"
DEFAULT_NETWORK="scrapper-client-network"
ENV_VARS=()


if [ "$#" -lt 2 ]; then
  echo -e "\e[31mUso: $0 -n <nombre_contenedor> -i <ip> [-p <puerto_host:puerto_contenedor>] [-e <env_var=valor>] [-network <red>] [-image <imagen>]\e[0m"
  exit 1
else
  echo -e "\e[32mIniciando script add-client.sh\e[0m"
fi

while [[ "$#" -gt 0 ]]; do
  case $1 in
    -n|--name) CONTAINER_NAME="$2"; shift 2 ;;
    -i|--ip) IP="$2"; shift 2 ;;
    -p|--port) PORT="$2"; shift 2 ;;
    -e|--env) ENV_VARS+=("$2"); shift 2 ;;
    -network|--network) NETWORK="$2"; shift 2 ;;
    -image|--image) IMAGE="$2"; shift 2 ;;
    *) echo -e "\e[31mOpción desconocida: $1\e[0m"; exit 1 ;;
  esac
done

if [ -z "$CONTAINER_NAME" ] || [ -z "$IP" ]; then
  echo -e "\e[31mError: Los parámetros -n y -i son obligatorios.\e[0m"
  exit 1
fi

PORT="${PORT:-$DEFAULT_PORT}"
NETWORK="${NETWORK:-$DEFAULT_NETWORK}"
IMAGE="${IMAGE:-$DEFAULT_IMAGE}"
VOLUMES=("${DEFAULT_VOLUMES[@]}")

if ! docker image inspect "$IMAGE" > /dev/null 2>&1; then
  echo -e "\e[33mImagen $IMAGE no encontrada localmente. Construyendo la imagen...\e[0m"
  docker build -t "$IMAGE" ./client || { echo -e "\e[31mError al construir la imagen $IMAGE.\e[0m"; exit 1; }
fi

DOCKER_CMD="docker run -d --name $CONTAINER_NAME --network $NETWORK --ip $IP -p $PORT"

for ENV in "${ENV_VARS[@]}"; do
  DOCKER_CMD+=" -e $ENV"
done

for VOL in "${VOLUMES[@]}"; do
  DOCKER_CMD+=" -v $VOL"
done

DOCKER_CMD+=" --privileged --cap-add=NET_ADMIN $IMAGE"

echo -e "\e[32mEjecutando: $DOCKER_CMD\e[0m"
eval "$DOCKER_CMD"

if [ $? -eq 0 ]; then
  echo -e "\e[32mContenedor $CONTAINER_NAME creado exitosamente.\e[0m"
else
  echo -e "\e[31mError al crear el contenedor $CONTAINER_NAME.\e[0m"
fi
