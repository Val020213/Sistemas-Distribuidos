#!/bin/bash

# Guardar la imagen my-router-image en un archivo tar
# docker save -o my-router-image.tar my-router-image:latest
# cp my-router-image.tar ./images

# # Guardar la imagen my-backend-image en un archivo tar
# docker save -o my-backend-image.tar my-backend-image:latest
# cp my-backend-image.tar ./images

# # Guardar la imagen my-frontend1-image en un archivo tar
# docker save -o my-frontend1-image.tar my-frontend1-image:latest
# cp my-frontend1-image.tar ./images

# # Guardar la imagen my-frontend2-image en un archivo tar
# docker save -o my-frontend2-image.tar my-frontend2-image:latest
# cp my-frontend2-image.tar ./images

docker save -o mongo.tar mongo:latest