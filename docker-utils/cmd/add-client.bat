@echo off
setlocal enabledelayedexpansion

set "DEFAULT_IMAGE=scrapper-client-image:latest"
set "DEFAULT_PORT=3000:3000"
set "DEFAULT_NETWORK=scrapper-client-network"
set "ENV_VARS="

if "%~1"=="" (
    echo Uso: %~nx0 -n ^<nombre_contenedor^> -i ^<ip^> [-p ^<puerto_host:puerto_contenedor^>] [-e ^<env_var=valor^>] [-network ^<red^>] [-image ^<imagen^>]
    exit /b 1
) else (
    echo Iniciando script add-client.bat
)

:parse_args
if "%~1"=="" goto :done_parse_args
if "%~1"=="-n" (
    set "CONTAINER_NAME=%~2"
    shift
    shift
    goto :parse_args
)
if "%~1"=="-i" (
    set "IP=%~2"
    shift
    shift
    goto :parse_args
)
if "%~1"=="-p" (
    set "PORT=%~2"
    shift
    shift
    goto :parse_args
)
if "%~1"=="-e" (
    set "ENV_VARS=!ENV_VARS! -e %~2"
    shift
    shift
    goto :parse_args
)
if "%~1"=="-network" (
    set "NETWORK=%~2"
    shift
    shift
    goto :parse_args
)
if "%~1"=="-image" (
    set "IMAGE=%~2"
    shift
    shift
    goto :parse_args
)
echo Opción desconocida: %~1
exit /b 1

:done_parse_args

if "%CONTAINER_NAME%"=="" (
    echo Error: El parámetro -n es obligatorio.
    exit /b 1
)
if "%IP%"=="" (
    echo Error: El parámetro -i es obligatorio.
    exit /b 1
)

if "%PORT%"=="" set "PORT=%DEFAULT_PORT%"
if "%NETWORK%"=="" set "NETWORK=%DEFAULT_NETWORK%"
if "%IMAGE%"=="" set "IMAGE=%DEFAULT_IMAGE%"

docker image inspect "%IMAGE%" >nul 2>&1
if errorlevel 1 (
    echo Imagen %IMAGE% no encontrada localmente. Construyendo la imagen...
    docker build -t "%IMAGE%" ./client
    if errorlevel 1 (
        echo Error al construir la imagen %IMAGE%.
        exit /b 1
    )
)

set "DOCKER_CMD=docker run -d --name %CONTAINER_NAME% --network %NETWORK% --ip %IP% -p %PORT% %ENV_VARS% --privileged --cap-add=NET_ADMIN %IMAGE%"

echo Ejecutando: %DOCKER_CMD%
%DOCKER_CMD%

if errorlevel 0 (
    echo Contenedor %CONTAINER_NAME% creado exitosamente.
) else (
    echo Error al crear el contenedor %CONTAINER_NAME%.
)