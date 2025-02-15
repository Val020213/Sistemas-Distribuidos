#!/bin/bash

set -e

cd ./src || exit 1

chmod +x ./server/server.sh

# Set the number of new instances you want to create
INSTANCES=1

# Check if scrapper-backend-image:latest exists; if not, build it
IMG_ID=$(docker images -q scrapper-backend-image:latest 2>/dev/null || true)
if [ -z "$IMG_ID" ]; then
    echo "scrapper-backend-image:latest not found. Building image..."
    docker build -t scrapper-backend-image:latest ./server || { echo "Error: Failed to build scrapper-backend-image."; exit 1; }
else
    echo "scrapper-backend-image:latest already exists."
fi

# Ensure the external network exists; if not, create it
if ! docker network ls --filter name=scrapper-server-network | grep -q scrapper-server-network; then
    echo "Network scrapper-server-network not found. Creating network..."
    docker network create scrapper-server-network --subnet=10.0.10.0/24
else
    echo "Network scrapper-server-network already exists."
fi

# Determine the last used instance number from existing backend containers
LAST_INSTANCE=0
for name in $(docker ps -a --format "{{.Names}}" | grep -E '^backend[0-9]+$' || true); do
    num=${name#backend}
    if [[ "$num" =~ ^[0-9]+$ ]]; then
        if [ "$num" -gt "$LAST_INSTANCE" ]; then
            LAST_INSTANCE=$num
        fi
    fi
done

echo "Last instance found: $LAST_INSTANCE"

# Calculate starting and ending instance numbers for new containers
START_INSTANCE=$((LAST_INSTANCE + 1))
END_INSTANCE=$((START_INSTANCE + INSTANCES - 1))

echo "New instances will be created starting from $START_INSTANCE to $END_INSTANCE."

# Create new backend and mongo instances without removing existing ones
for ((i = START_INSTANCE; i <= END_INSTANCE; i++)); do
    backendPort=$((8080 + i))
    mongoPort=$((27017 + i))
    echo "Creating mongo instance $i..."
    docker run -d --restart unless-stopped --name mongo_bp$i \
    -p "$mongoPort":27017 \
    -v mongo_volume_bp$i:/data/db \
    --network scrapper-server-network \
    --ip 10.0.10.2$i \
    --network-alias mongo$i \
    -e MONGO_INITDB_ROOT_USERNAME=melkey \
    -e MONGO_INITDB_ROOT_PASSWORD=password1234 \
    -e MONGO_INITDB_DATABASE=mongo_bp \
    --health-cmd "mongosh --eval 'db.adminCommand(\"ping\")'" \
    --health-interval 5s \
    --health-timeout 5s \
    --health-retries 3 \
    --health-start-period 60s \
    --ulimit nofile=64000:64000 \
    mongo:latest mongod --quiet --logpath /var/log/mongodb/mongod.log



    echo "Creating backend instance $i..."
    docker run -d --restart unless-stopped --name backend$i \
        -p "$backendPort":8080 \
        -v "$PWD/server":/app \
        --network scrapper-server-network \
        --ip 10.0.10.1$i \
        --cap-add NET_ADMIN \
        --privileged \
        -e BLUEPRINT_DB_HOST=mongo$i \
        -e BLUEPRINT_DB_PORT=27017 \
        -e BLUEPRINT_DB_USER=melkey \
        -e BLUEPRINT_DB_PASSWORD=password1234 \
        -e BLUEPRINT_DB_NAME=mongo_bp \
        -e PORT=8080 \
        -e TASK_QUEUE_SIZE=10 \
        -e NUM_WORKERS=5 \
        scrapper-backend-image:latest
done

echo "All new instances created."
