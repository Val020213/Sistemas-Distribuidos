#!/bin/bash
cd ./src || exit 1

# Verificar si las imágenes existen
if ! docker images | grep -q "scrapper-router-image"; then
    echo "Construyendo la imagen scrapper-router-image..."
    docker-compose build router
fi

if ! docker images | grep -q "my-backend-image"; then
    echo "Construyendo la imagen my-backend-image..."
    docker-compose build backend
fi

if ! docker images | grep -q "scrapper-client-image"; then
    echo "Construyendo la imagen scrapper-client-image..."
    docker-compose build frontend1
fi

# Levantar los contenedores sin volver a construir las imágenes
docker-compose up

sleep 10 # Delay para esperar a que los contenedores estén listos
xdg-open http://localhost:3000 # Abrir el navegador para el frontend1
xdg-open http://localhost:3001 # Abrir el navegador para el frontend2