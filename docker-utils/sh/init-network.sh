#!/bin/bash

set -e

cd ./src || exit 1

# Ensure scrapper-server-network exists
if ! docker network inspect scrapper-server-network >/dev/null 2>&1; then
    echo "Network scrapper-server-network not found. Creating network..."
    docker network create scrapper-server-network --subnet=10.0.10.0/24
else
    echo "Network scrapper-server-network exists."
fi

# Ensure scrapper-client-network exists
if ! docker network inspect scrapper-client-network >/dev/null 2>&1; then
    echo "Network scrapper-client-network not found. Creating network..."
    docker network create scrapper-client-network --subnet=10.0.11.0/24
else
    echo "Network scrapper-client-network exists."
fi

# Ensure scrapper-router-image exists
if ! docker image inspect scrapper-router-image:latest >/dev/null 2>&1; then
    echo "Image scrapper-router-image not found. Building image..."
    docker build -t scrapper-router-image:latest ./router
    echo "Image scrapper-router-image created."
else
    echo "Image scrapper-router-image exists."
fi

# Ensure router container is not running
if docker container inspect router >/dev/null 2>&1; then
    docker container stop router
    docker container rm router
    echo "Container router removed."
fi

# Run the router container
docker run -d --rm --name router --cap-add NET_ADMIN -e PYTHONUNBUFFERED=1 scrapper-router-image

# Connect router to both networks
docker network connect --ip 10.0.10.254 scrapper-server-network router
docker network connect --ip 10.0.11.254 scrapper-client-network router

echo "Container router executed and connected to networks."
