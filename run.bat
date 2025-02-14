@echo off

REM Navegar al directorio que contiene el docker-compose.yml
cd src || exit /b 1

@echo off
REM Verificar si las imágenes existen
docker images | findstr "scrapper-router-image" >nul
if %ERRORLEVEL% NEQ 0 (
    echo "Construyendo la imagen scrapper-router-image..."
    docker-compose build router
)

docker images | findstr "scrapper-server-image" >nul
if %ERRORLEVEL% NEQ 0 (
    echo "Construyendo la imagen scrapper-server-image..."
    docker-compose build backend
)

docker images | findstr "scrapper-client-image" >nul
if %ERRORLEVEL% NEQ 0 (
    echo "Construyendo la imagen scrapper-client-image..."
    docker-compose build frontend
)


REM Levantar los contenedores sin volver a construir las imágenes
docker-compose up


REM Dar tiempo para que los contenedores se levanten
timeout /t 10

REM Abrir el navegador para el frontend1
start http://localhost:3000

REM Abrir el navegador para el frontend2
start http://localhost:3001