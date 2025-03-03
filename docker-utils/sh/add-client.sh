#!/bin/bash

set -e

cd ./src || exit 1

# Número de nuevas instancias a crear
INSTANCES=1

# Revisar si la imagen de frontend existe, si no, la construye
IMG_ID=$(docker images -q scrapper-client-image:latest 2>/dev/null || true)
if [ -z "$IMG_ID" ]; then
    echo "scrapper-client-image:latest no encontrado. Construyendo imagen..."
    docker build -t scrapper-client-image:latest ./client || { echo "Error: No se pudo construir la imagen de frontend."; exit 1; }
else
    echo "scrapper-client-image:latest ya existe."
fi

# Verificar si la red "scrapper-client-network" existe; si no, crearla
if ! docker network ls --filter name=scrapper-client-network | grep -q scrapper-client-network; then
    echo "Red scrapper-client-network no encontrada. Creándola..."
    docker network create scrapper-client-network --subnet=10.0.11.0/24
else
    echo "Red scrapper-client-network ya existe."
fi

# Determinar la última instancia de frontend creada
LAST_INSTANCE=0
for name in $(docker ps -a --format "{{.Names}}" | grep -E '^frontend[0-9]+$' || true); do
    num=${name#frontend} # Extraer el número después de "frontend"
    if [[ "$num" =~ ^[0-9]+$ ]]; then
        if [ "$num" -gt "$LAST_INSTANCE" ]; then
            LAST_INSTANCE=$num
        fi
    fi
done

echo "Última instancia encontrada: $LAST_INSTANCE"

# Calcular las nuevas instancias a crear
START_INSTANCE=$((LAST_INSTANCE + 1))
END_INSTANCE=$((START_INSTANCE + INSTANCES - 1))

echo "Creando nuevas instancias desde $START_INSTANCE hasta $END_INSTANCE."

# Crear nuevas instancias de frontend
for ((i = START_INSTANCE; i <= END_INSTANCE; i++)); do
    frontendPort=$((3000 + i))
    echo "Creando frontend instance $i..."

    docker run -d --restart unless-stopped --name frontend$i \
        -p "$frontendPort":3000 \
        -v "$PWD/client":/app \
        -v /app/node_modules \
        -v /app/.next \
        --network scrapper-client-network \
        --ip 10.0.11.1$i \
        --cap-add NET_ADMIN \
        --privileged \
        --env-file "$PWD/client/.env" \
        -e NEXT_PUBLIC_API_URL=http://10.0.11.254:8080 \
        scrapper-client-image:latest

done

echo "Todas las nuevas instancias han sido creadas."
