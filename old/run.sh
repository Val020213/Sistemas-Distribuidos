#!/bin/bash
cd ./src || exit 1

# Verificar si las imágenes existen
if ! docker images | grep -q "scrapper-router-image"; then
    echo "Construyendo la imagen del router..."
    docker-compose build router
fi

# if ! docker images | grep -q "my-backend-image"; then
#     echo "Construyendo la imagen my-backend-image..."
#     docker-compose build backend
# fi

if ! docker images | grep -q "scrapper-client-image"; then
    echo "Construyendo la imagen del cliente..."
    docker-compose build frontend
fi

# if ! docker images | grep -q "my-frontend2-image"; then
#     echo "Construyendo la imagen my-frontend2-image..."
#     docker-compose build frontend2
# fi

# Levantar los contenedores sin volver a construir las imágenes
docker-compose up

sleep 10 # Delay para esperar a que los contenedores estén listos
xdg-open http://localhost:3000 # Abrir el navegador para el frontend1
# xdg-open http://localhost:3001 # Abrir el navegador para el frontend2