#!/bin/bash
cd ./images || exit 1

# Cargar las imágenes Docker desde los archivos tar
docker load -i ./images/my-router-image.tar
docker load -i ./images/my-backend-image.tar
docker load -i ./images/my-frontend1-image.tar
docker load -i ./images/my-frontend2-image.tar