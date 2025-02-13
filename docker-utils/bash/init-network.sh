#!/bin/bash
cd ./src || exit 1

# Configuración
export COMPOSE_PROJECT_NAME=scrapper
export ROUTER_IMAGE="scrapper-router-image"
export ROUTER_CONTAINER="scrapper-router"

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Función para encontrar redes conflictivas
find_conflicting_networks() {
    TARGET_SUBNET=$1
    docker network ls --format '{{.ID}}' | while read -r NET_ID; do
        SUBNET=$(docker network inspect $NET_ID --format '{{range .IPAM.Config}}{{.Subnet}}{{end}}')
        if [ "$SUBNET" = "$TARGET_SUBNET" ]; then
            echo -e "${RED}Red conflictiva encontrada:${NC}"
            docker network inspect $NET_ID --format 'Nombre: {{.Name}}\nID: {{.ID}}\nSubnet: {{range .IPAM.Config}}{{.Subnet}}{{end}}\n'
            echo "Eliminando red conflictiva..."
            docker network rm $NET_ID
        fi
    done
}

# Limpiar redes existentes
echo -e "${GREEN}Buscando redes conflictivas...${NC}"
find_conflicting_networks "10.0.10.0/24"
find_conflicting_networks "10.0.11.0/24"
#!/bin/bash
cd ../src || exit 1

# Configuración
export COMPOSE_PROJECT_NAME=scrapper
export ROUTER_IMAGE="scrapper-router-image"
export ROUTER_CONTAINER="scrapper-router"

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Función para verificar y crear el contenedor
create_router() {
    echo -e "${GREEN}Verificando existencia del router...${NC}"
    
    if ! docker inspect $ROUTER_CONTAINER &> /dev/null; then
        echo -e "${RED}Router no encontrado, creando nuevo...${NC}"
        
        # Eliminar residuos previos
        docker rm -f $ROUTER_CONTAINER 2>/dev/null || true
        
        # Construir imagen si es necesario
        if ! docker image inspect $ROUTER_IMAGE:latest &> /dev/null; then
            echo -e "${GREEN}Construyendo imagen del router...${NC}"
            docker-compose build --no-cache router
        fi
        
        # Crear contenedor
        echo -e "${GREEN}Iniciando nuevo router...${NC}"
        docker-compose up -d --force-recreate router
        
        # Esperar inicialización
        sleep 5
        echo -e "${GREEN}Verificando estado...${NC}"
        docker ps -f "name=$ROUTER_CONTAINER"
    else
        echo -e "${GREEN}Router ya existe, reiniciando...${NC}"
        docker restart $ROUTER_CONTAINER
    fi
}

# Ejecutar flujo principal
create_router

# Verificación final
echo -e "\n${GREEN}Estado del router:${NC}"
docker inspect $ROUTER_CONTAINER --format '
Nombre: {{.Name}}
Estado: {{.State.Status}}
IPs: {{range .NetworkSettings.Networks}}{{.IPAddress}} {{end}}
Puertos: {{range $p, $conf := .NetworkSettings.Ports}}{{$p}} {{end}}'

echo -e "\n${GREEN}Prueba de conectividad interna:${NC}"
docker exec $ROUTER_CONTAINER ping -c 2 10.0.10.254
docker exec $ROUTER_CONTAINER ping -c 2 10.0.11.254

docker exec scrapper-router sh -c '
iptables -t nat -A POSTROUTING -s 10.0.10.0/24 -o eth0 -j MASQUERADE;
iptables -t nat -A POSTROUTING -s 10.0.11.0/24 -o eth1 -j MASQUERADE;
iptables -A FORWARD -i eth1 -o eth0 -j ACCEPT;
iptables -A FORWARD -i eth0 -o eth1 -j ACCEPT'

docker exec scrapper-router ip route replace default via 10.0.10.254 dev eth1
docker exec scrapper-router ip route add 10.0.11.0/24 dev eth0
docker exec scrapper-router ip route add 10.0.10.0/24 dev eth1

docker exec scrapper-router iptables -P FORWARD ACCEPT
