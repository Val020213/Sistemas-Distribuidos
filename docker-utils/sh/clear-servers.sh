#!/bin/bash

if [ "$#" -eq 0 ]; then
    echo "Usage: $0 <number1> <number2> ..."
    exit 1
fi

for i in "$@"; do
    echo "Stopping container: backend${i}"
    docker stop "backend${i}"
    
    echo "Stopping container: mongo_bp${i}"
    docker stop "mongo_bp${i}"

done

echo "Pruning stopped containers and unused volumes..."
docker container prune -f

for i in "$@"; do
  
    echo "Deleting volume: mongo_volume_bp${i}"
    docker volume rm "mongo_volume_bp${i}"
done


docker volume prune -f

echo "Done."