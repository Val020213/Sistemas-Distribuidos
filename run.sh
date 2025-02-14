#!/bin/bash
cd ./src || exit 1

chmod +x ./scrapper_server/server.sh


# Verificar si las imágenes existen
if ! docker images | grep -q "scrapper-router-image"; then
    echo "Construyendo la imagen scrapper-router-image..."
    docker-compose build router
fi

if ! docker images | grep -q "scrapper-server-image"; then
    echo "Construyendo la imagen scrapper-server-image..."
    docker-compose build backend
fi

if ! docker images | grep -q "scrapper-client-image"; then
    echo "Construyendo la imagen scrapper-client-image..."
    docker-compose build frontend
fi

# Levantar los contenedores sin volver a construir las imágenes
docker-compose up

