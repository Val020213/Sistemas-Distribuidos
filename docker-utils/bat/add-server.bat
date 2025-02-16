@echo off
setlocal enabledelayedexpansion


cd ./src || exit 1

:: Set the number of new instances you want to create
set INSTANCES=1

:: Check if scrapper-backend-image:latest exists; if not, build it
for /f "delims=" %%a in ('docker images -q scrapper-backend-image:latest 2^>nul') do (
  set IMG_ID=%%a
)
if "%IMG_ID%"=="" (
    echo scrapper-backend-image:latest not found. Building image...
    docker build -t scrapper-backend-image:latest ./server
    if errorlevel 1 (
        echo Error: Failed to build scrapper-backend-image.
        exit /b 1
    )
) else (
    echo scrapper-backend-image:latest already exists.
)

:: Ensure the external network exists; if not, create it
docker network ls --filter name=scrapper-server-network | findstr /i scrapper-server-network >nul
if errorlevel 1 (
    echo Network scrapper-server-network not found. Creating network...
    docker network create scrapper-server-network --subnet=10.0.10.0/24
) else (
    echo Network scrapper-server-network already exists.
)

:: Determine the last used instance number from existing backend containers
set LAST_INSTANCE=0
for /f "tokens=*" %%a in ('docker ps -a --format "{{.Names}}" ^| findstr /i "^backend"') do (
    set "name=%%a"
    set "num=!name:~7!"
    if not "!num!"=="" (
        if !num! gtr !LAST_INSTANCE! (
            set LAST_INSTANCE=!num!
        )
    )
)
echo Last instance found: %LAST_INSTANCE%

:: Calculate starting and ending instance numbers for new containers
set /a START_INSTANCE=LAST_INSTANCE+1
set /a END_INSTANCE=START_INSTANCE+INSTANCES-1
echo New instances will be created starting from %START_INSTANCE% to %END_INSTANCE%.

:: Create new backend and mongo instances without removing existing ones
for /L %%i in (%START_INSTANCE%,1,%END_INSTANCE%) do (
    set /a backendPort=8080+%%i
    set /a mongoPort=27017+%%i
    set /a grpcPort=50050+%%i
    echo Creating mongo instance %%i...
    docker run -d --restart unless-stopped --name mongo_bp%%i -p !mongoPort!:27017 -v mongo_volume_bp%%i:/data/db --network scrapper-server-network --ip 10.0.10.2%%i --network-alias mongo%%i -e MONGO_INITDB_ROOT_USERNAME=melkey -e MONGO_INITDB_ROOT_PASSWORD=password1234 -e MONGO_INITDB_DATABASE=mongo_bp --health-cmd "mongosh --eval ""db.adminCommand('ping')""" --health-interval 5s --health-timeout 5s --health-retries 3 --health-start-period 15s mongo:latest mongod --quiet --logpath /dev/null

    echo Creating backend instance %%i...
    docker run -d --restart unless-stopped --name backend%%i -p !backendPort!:8080 -p !grpcPort!:50051 -v "%cd%\server":/app  -v "go-mod-cache:/go/pkg/mod" --network scrapper-server-network --ip 10.0.10.1%%i --cap-add NET_ADMIN --privileged -e BLUEPRINT_DB_HOST=mongo%%i -e BLUEPRINT_DB_PORT=27017 -e BLUEPRINT_DB_USER=melkey -e BLUEPRINT_DB_PASSWORD=password1234 -e BLUEPRINT_DB_NAME=mongo_bp -e PORT=8080 -e RPC_PORT=50051 -e IP_ADDRESS=10.0.10.1%%i -e TASK_QUEUE_SIZE=10 -e NUM_WORKERS=5 scrapper-backend-image:latest

)

echo All new instances created.
pause
