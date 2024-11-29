#!/bin/bash
docker-compose down -v # en caso de que ya exista un contenedor con el mismo nombre
docker-compose up --build

sleep 10 # Delay para esperar a que los contenedores estén listos
xdg-open http://localhost:3000 # Abrir el navegador para el frontend1
xdg-open http://localhost:3001 # Abrir el navegador para el frontend2