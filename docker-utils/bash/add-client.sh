#!/bin/bash
set -e

# Configuración de rutas
CLIENT_DIR="$(pwd)/src/client"
DOCKERFILE_PATH="$CLIENT_DIR/Dockerfile"
ENV_FILE="$CLIENT_DIR/.env"
IMAGE_NAME="scrapper-client-image:latest"
NETWORK="scrapper-client-network"
BASE_IP="10.0.11"
API_ENDPOINT="10.0.10.2:8080"

# Verificar Dockerfile
[ -f "$DOCKERFILE_PATH" ] || { echo "Error: Dockerfile no encontrado"; exit 1; }

# Construir imagen si no existe
if ! docker image inspect "$IMAGE_NAME" &>/dev/null; then
    echo "Construyendo imagen del cliente..."
    docker build -t "$IMAGE_NAME" -f "$DOCKERFILE_PATH" "$CLIENT_DIR"
fi

# Manejo de .env
ENV_FILE_OPTION=""
if [ -f "$ENV_FILE" ]; then
    ENV_FILE_OPTION="--env-file $ENV_FILE"
else
    read -p "¿Continuar sin .env? (y/n) " -n 1 -r
    [[ $REPLY =~ ^[Yy]$ ]] || exit 1
fi

# Calcular siguiente IP
LAST_IP=$(docker network inspect "$NETWORK" --format '{{range .Containers}}{{.IPv4Address}}{{end}}' 2>/dev/null | \
          grep -oP '\d+\.\d+\.\d+\.\d+' | sort -t. -k4 -n | tail -1 || echo "10.0.11.1")

LAST_OCTET=${LAST_IP##*.}
NEXT_OCTET=$((LAST_OCTET + 1))

# Corregir cálculo de IP (evitar 254 y 255)
if [[ $NEXT_OCTET -ge 254 ]]; then
    NEXT_OCTET=2
fi
NEW_IP="${BASE_IP}.${NEXT_OCTET}"

# Calcular puerto válido (3000-65535)
LAST_PORT=$(docker ps --format '{{.Ports}}' | grep -oP '\d+(?=->3000)' | sort -n | tail -1 || echo 2999)
NEXT_PORT=$((LAST_PORT + 1))

# Asegurar puerto mínimo de 3000
if [[ $NEXT_PORT -lt 3000 ]]; then
    NEXT_PORT=3000
fi

# Verificar puerto disponible
while true; do
    if ! ss -tuln | grep -q ":${NEXT_PORT} "; then
        break
    fi
    ((NEXT_PORT++))
    
    if [[ $NEXT_PORT -gt 65535 ]]; then
        echo "Error: No hay puertos disponibles"
        exit 1
    fi
done

# Crear contenedor
echo "Creando cliente:"
echo "Nombre: frontend-${NEXT_OCTET}"
echo "IP: ${NEW_IP}"
echo "Puerto: ${NEXT_PORT}"

docker run -d \
  --name "frontend-${NEXT_OCTET}" \
  --network "$NETWORK" \
  --ip "$NEW_IP" \
  --restart unless-stopped \
  --cap-add NET_ADMIN \
  --privileged \
  -v "${CLIENT_DIR}:/app" \
  -v "/app/node_modules" \
  -v "/app/.next" \
  -p "${NEXT_PORT}:3000" \
  -e NEXT_PUBLIC_API_URL="http://${API_ENDPOINT}" \
  $ENV_FILE_OPTION \
  "$IMAGE_NAME"

echo "Cliente creado exitosamente!"

# Ver IPs asignadas
docker network inspect scrapper-client-network --format '{{range .Containers}}{{.Name}} - {{.IPv4Address}}{{end}}'

# Ver puertos en uso
ss -tuln | grep 'LISTEN'