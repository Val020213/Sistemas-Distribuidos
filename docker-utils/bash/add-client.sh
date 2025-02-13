#!/bin/bash
set -e
cd ./src/ || exit 1
# Configuration
CLIENT_DIR="$(pwd)/client"
COMPOSE_FILE="docker-compose-client.yaml"
# NETWORK="scrapper-client-network"
# BASE_IP="10.0.11"
# API_IP="10.0.10.2"
# API_PORT="8080"
# SUBNET="${BASE_IP}.0/24"

# Validate essential files
[ -f "$COMPOSE_FILE" ] || { echo "Error: $COMPOSE_FILE not found"; exit 1; }
[ -d "$CLIENT_DIR" ] || { echo "Error: client directory missing"; exit 1; }

# # Create network if not exists
# if ! docker network inspect "$NETWORK" &>/dev/null; then
#     echo "Creating network $NETWORK..."
#     docker network create \
#         --driver bridge \
#         --subnet "$SUBNET" \
#         --gateway "${BASE_IP}.254" \
#         "$NETWORK"
# fi

# # Calculate next client number
# LAST_CLIENT=$(docker ps -a --filter "name=frontend" --format "{{.Names}}" | 
#     grep -oP '\d+$' | sort -n | tail -n1 || echo "0")
# NEXT_CLIENT=$((LAST_CLIENT + 1))

# # Calculate next available IP
# IP_LIST=$(docker network inspect "$NETWORK" --format '{{range .Containers}}{{.IPv4Address}} {{end}}')
# VALID_IPS=$(echo "$IP_LIST" | grep -oP '\d+\.\d+\.\d+\.\d+' | 
#     grep -vE '\.254$|\.255$' | sort -t. -k4n)

# if [ -z "$VALID_IPS" ]; then
#     NEXT_OCTET=2  # Start from .2
# else
#     LAST_IP=$(echo "$VALID_IPS" | tail -n1)
#     LAST_OCTET=${LAST_IP##*.}
#     NEXT_OCTET=$((LAST_OCTET + 1))
    
#     # Reset sequence if approaching reserved IPs
#     if [ "$NEXT_OCTET" -ge 254 ]; then
#         echo "Warning: IP range exhausted, recycling from .2"
#         NEXT_OCTET=2
#     fi
# fi

# NEW_IP="${BASE_IP}.${NEXT_OCTET}"

# # Find available port
# PORT_START=3000
# PORT_END=4000
# NEXT_PORT=$(comm -23 \
#     <(seq $PORT_START $PORT_END | sort) \
#     <(docker ps --format '{{.Ports}}' | 
#         grep -oP '\d+(?=->3000)' | sort -n) |
#     head -n1)

# # Generate dynamic compose override
# OVERRIDE_FILE=".client-override.yml"
# cat > "$OVERRIDE_FILE" <<EOF
# version: '3.8'

# services:
#   frontend:
#     image: scrapper-client-image:latest
#     container_name: frontend${NEXT_CLIENT}
#     networks:
#       $NETWORK:
#         ipv4_address: $NEW_IP
#     ports:
#       - "$NEXT_PORT:3000"
#     environment:
#       - NEXT_PUBLIC_API_URL=http://$API_IP:$API_PORT

# networks:
#   $NETWORK:
#     external: true
# EOF

# Run docker-compose
# echo "Creating client ${NEXT_CLIENT}:"
# echo "IP: $NEW_IP"
# echo "Port: $NEXT_PORT"
# echo "Container: frontend${NEXT_CLIENT}"

docker-compose -f "$COMPOSE_FILE" up 

# # Cleanup
# rm "$OVERRIDE_FILE"

echo "Client created successfully!"