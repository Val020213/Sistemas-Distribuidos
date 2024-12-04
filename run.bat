@echo off

REM Navegar al directorio que contiene el docker-compose.yml
cd src || exit /b 1

@echo off
REM Verificar si las imágenes existen
docker images | findstr "my-router-image" >nul
if %ERRORLEVEL% NEQ 0 (
    echo "Construyendo la imagen my-router-image..."
    docker-compose build router
)

docker images | findstr "my-backend-image" >nul
if %ERRORLEVEL% NEQ 0 (
    echo "Construyendo la imagen my-backend-image..."
    docker-compose build backend
)

docker images | findstr "my-frontend1-image" >nul
if %ERRORLEVEL% NEQ 0 (
    echo "Construyendo la imagen my-frontend1-image..."
    docker-compose build frontend1
)

docker images | findstr "my-frontend2-image" >nul
if %ERRORLEVEL% NEQ 0 (
    echo "Construyendo la imagen my-frontend2-image..."
    docker-compose build frontend2
)

REM Levantar los contenedores sin volver a construir las imágenes
docker-compose up


REM Dar tiempo para que los contenedores se levanten
timeout /t 10

REM Abrir el navegador para el frontend1
start http://localhost:3000

REM Abrir el navegador para el frontend2
start http://localhost:3001