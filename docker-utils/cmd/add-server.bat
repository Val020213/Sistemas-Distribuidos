@echo off

REM Enable ANSI escape codes
for /F "tokens=2 delims==" %%i in ('"prompt $E" ^& for %%i in (1) do rem') do set "ESC=%%i"

REM ANSI color codes
set GREEN=%ESC%[32m
set NC=%ESC%[0m

REM Default values
set DEFAULT_NETWORK_NAME=scrapper-server-network
set DEFAULT_SUBNET=10.0.0.0/24
set DEFAULT_MONGO_IMAGE=mongo:6.0
set DEFAULT_SERVER_IMAGE=scrapper_server

REM Default values for MongoDB
set MONGO_VOLUME=mongo_data
set MONGO_PORT=27018
set MONGO_USERNAME=admin
set MONGO_PASSWORD=password

REM Read parameters or use default values
:parse_args
if "%1"=="" goto end_args
if "%1"=="-m" set MONGO_CONTAINER_NAME=%2
if "%1"=="-i" set MONGO_IP=%2
if "%1"=="-s" set SERVER_CONTAINER_NAME=%2
if "%1"=="-j" set SERVER_IP=%2
if "%1"=="-p" set SERVER_PORT=%2
shift
shift
goto parse_args
:end_args

REM Validate required parameters
if "%MONGO_CONTAINER_NAME%"=="" (
    echo Missing required parameters. Usage: %0 -m ^<mongo_container_name^> -i ^<mongo_ip^> -s ^<server_container_name^> -j ^<server_ip^> -p ^<server_port^>
    exit /b 1
)
if "%MONGO_IP%"=="" (
    echo Missing required parameters. Usage: %0 -m ^<mongo_container_name^> -i ^<mongo_ip^> -s ^<server_container_name^> -j ^<server_ip^> -p ^<server_port^>
    exit /b 1
)
if "%SERVER_CONTAINER_NAME%"=="" (
    echo Missing required parameters. Usage: %0 -m ^<mongo_container_name^> -i ^<mongo_ip^> -s ^<server_container_name^> -j ^<server_ip^> -p ^<server_port^>
    exit /b 1
)
if "%SERVER_IP%"=="" (
    echo Missing required parameters. Usage: %0 -m ^<mongo_container_name^> -i ^<mongo_ip^> -s ^<server_container_name^> -j ^<server_ip^> -p ^<server_port^>
    exit /b 1
)
if "%SERVER_PORT%"=="" (
    echo Missing required parameters. Usage: %0 -m ^<mongo_container_name^> -i ^<mongo_ip^> -s ^<server_container_name^> -j ^<server_ip^> -p ^<server_port^>
    exit /b 1
)

set NETWORK_NAME=%DEFAULT_NETWORK_NAME%
set MONGO_CONTAINER_NAME=%MONGO_CONTAINER_NAME%
set MONGO_IP=%MONGO_IP%
set SERVER_CONTAINER_NAME=%SERVER_CONTAINER_NAME%
set SERVER_IP=%SERVER_IP%
set SERVER_PORT=%SERVER_PORT%

REM Create volume for MongoDB if it doesn't exist
docker volume create %MONGO_VOLUME%
docker run -d ^
    --name %MONGO_CONTAINER_NAME% ^
    --rm ^
    --network %NETWORK_NAME% ^
    --ip %MONGO_IP% ^
    -p %MONGO_PORT%:%MONGO_PORT% ^
    -e MONGO_INITDB_ROOT_USERNAME=%MONGO_USERNAME% ^
    -e MONGO_INITDB_ROOT_PASSWORD=%MONGO_PASSWORD% ^
    -v %MONGO_VOLUME%:/data/db ^
    %DEFAULT_MONGO_IMAGE%

echo %GREEN%MongoDB has been successfully started in container '%MONGO_CONTAINER_NAME%' with IP address %MONGO_IP%.%NC%

REM Build the server image from the Dockerfile
docker image inspect "%DEFAULT_SERVER_IMAGE%" >nul 2>&1
if errorlevel 1 (
    echo %GREEN%Docker image '%DEFAULT_SERVER_IMAGE%' not found locally. Building the image now...%NC%
    docker build -t "%DEFAULT_SERVER_IMAGE%" .\src\scrapper_server
    if errorlevel 1 (
        echo %GREEN%An error occurred while building the Docker image '%DEFAULT_SERVER_IMAGE%'.%NC%
        exit /b 1
    )
)

REM Start the server in the custom network with a fixed IP
docker run -d ^
    --name %SERVER_CONTAINER_NAME% ^
    --rm ^
    --network %NETWORK_NAME% ^
    --ip %SERVER_IP% ^
    -p %SERVER_PORT%:8080 ^
    -e PORT=%SERVER_PORT% ^
    -e BLUEPRINT_DB_HOST=%MONGO_IP% ^
    -e BLUEPRINT_DB_PORT=%MONGO_PORT% ^
    -e BLUEPRINT_DB_USER=%MONGO_USERNAME% ^
    -e BLUEPRINT_DB_PASSWORD=%MONGO_PASSWORD% ^
    -v %cd%:/app ^
    -v /app/go/pkg ^
    %DEFAULT_SERVER_IMAGE%

echo %GREEN%Server has been successfully started in container '%SERVER_CONTAINER_NAME%', listening on port %SERVER_PORT% with IP address %SERVER_IP%.%NC%