cd .\src\ || exit /b 1

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


docker images | find "%ROUTER_IMAGE%" >nul 2>&1 || (
    echo Building the Docker image %ROUTER_IMAGE%...
    docker build -t %ROUTER_IMAGE% .\router
)

echo Starting the network and router...
docker network ls | find "%SERVER_NETWORK%" && (
    echo Disconnecting existing containers from %SERVER_NETWORK%...
    for /f "tokens=*" %%i in ('docker network inspect "%SERVER_NETWORK%" -f "{{range .Containers}}{{.Name}} {{end}}"') do docker network disconnect "%SERVER_NETWORK%" %%i
    echo Removing existing network %SERVER_NETWORK%...
    docker network rm "%SERVER_NETWORK%"
)
echo Creating network %SERVER_NETWORK% with subnet %ROUTER_SUBNET_SERVER%...
docker network create --subnet="%ROUTER_SUBNET_SERVER%" "%SERVER_NETWORK%"

echo Creating network %CLIENT_NETWORK% with subnet %ROUTER_SUBNET_CLIENT%...
docker network create --subnet="%ROUTER_SUBNET_CLIENT%" "%CLIENT_NETWORK%"


echo Removing existing router container %ROUTER_CONTAINER% if it exists...
docker rm -f %ROUTER_CONTAINER% 2>nul || (exit /b 0)
if not exist %ROUTER_VOLUME% (
    echo Creating router volume directory...
    mkdir %ROUTER_VOLUME%
)

echo Running router container %ROUTER_CONTAINER%...

docker run -d --name %ROUTER_CONTAINER% --hostname %ROUTER_HOSTNAME% ^
    --network %SERVER_NETWORK% --ip %ROUTER_IP_SERVER% ^
    --cap-add=NET_ADMIN --privileged ^
    --sysctl net.ipv4.ip_forward=1 ^
    -v %ROUTER_VOLUME% %ROUTER_IMAGE%:latest

echo Connecting router container %ROUTER_CONTAINER% to client network %CLIENT_NETWORK% with IP %ROUTER_IP_CLIENT%...
docker network connect --ip %ROUTER_IP_CLIENT% %CLIENT_NETWORK% %ROUTER_CONTAINER%
