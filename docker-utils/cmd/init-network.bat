cd ./src || exit /b 1

set CLIENT_NETWORK=scrapper-client-network
set SERVER_NETWORK=scrapper-server-network
set ROUTER_IMAGE=scrapper-router-image
set ROUTER_CONTAINER=scrapper-router
set ROUTER_HOSTNAME=scrapper-router
set ROUTER_ALIAS=scrapper-router
set ROUTER_SUBNET_SERVER=10.0.10.0/24
set ROUTER_SUBNET_CLIENT=10.0.11.0/24
set ROUTER_IP_SERVER=10.0.10.254
set ROUTER_IP_CLIENT=10.0.11.254
set ROUTER_VOLUME=./router:/app

set GREEN=[32m
set NC=[0m

docker images | findstr "%ROUTER_IMAGE%" || (
    echo %GREEN%Construyendo la imagen %ROUTER_IMAGE%...%NC%
    docker-compose build router
)

echo %GREEN%Iniciando la red y router...%NC%
docker network ls | findstr "%SERVER_NETWORK%" && (
    for /f "tokens=*" %%i in ('docker network inspect "%SERVER_NETWORK%" -f "{{range .Containers}}{{.Name}} {{end}}"') do docker network disconnect "%SERVER_NETWORK%" %%i
    docker network rm "%SERVER_NETWORK%"
)
docker network create --subnet="%ROUTER_SUBNET_SERVER%" "%SERVER_NETWORK%"

docker rm -f %ROUTER_CONTAINER% 2>nul || (exit /b 0)

docker run -d --name %ROUTER_CONTAINER% --hostname %ROUTER_HOSTNAME% ^
    --network %SERVER_NETWORK% --ip %ROUTER_IP_SERVER% ^
    --cap-add=NET_ADMIN --privileged ^
    --sysctl net.ipv4.ip_forward=1 ^
    -v %ROUTER_VOLUME% %ROUTER_IMAGE%:latest

docker network connect --ip %ROUTER_IP_CLIENT% %CLIENT_NETWORK% %ROUTER_CONTAINER%