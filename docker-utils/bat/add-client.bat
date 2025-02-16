@echo off
setlocal enabledelayedexpansion

cd ./src || exit 1

:: Número de nuevas instancias a crear
set INSTANCES=1

:: Revisar si la imagen de frontend existe, si no, la construye
for /f "delims=" %%a in ('docker images -q scrapper-client-image:latest 2^>nul') do set IMG_ID=%%a
if "%IMG_ID%"=="" (
    echo scrapper-client-image:latest no encontrado. Construyendo imagen...
    docker build -t scrapper-client-image:latest ./client
    if errorlevel 1 (
        echo Error: No se pudo construir la imagen de frontend.
        exit /b 1
    )
) else (
    echo scrapper-client-image:latest ya existe.
)

:: Verificar si la red "scrapper-client-network" existe; si no, crearla
docker network ls --filter name=scrapper-client-network | findstr /i scrapper-client-network >nul
if errorlevel 1 (
    echo Red scrapper-client-network no encontrada. Creándola...
    docker network create scrapper-client-network --subnet=10.0.11.0/24
) else (
    echo Red scrapper-client-network ya existe.
)

:: Determinar la última instancia de frontend creada
set LAST_INSTANCE=0
for /f "tokens=*" %%a in ('docker ps -a --format "{{.Names}}" ^| findstr /i "^frontend"') do (
    set "name=%%a"
    set "num=!name:~8!"  :: Extraer el número después de "frontend"
    if not "!num!"=="" (
        if !num! gtr !LAST_INSTANCE! (
            set LAST_INSTANCE=!num!
        )
    )
)
echo Última instancia encontrada: %LAST_INSTANCE%

:: Calcular las nuevas instancias a crear
set /a START_INSTANCE=LAST_INSTANCE+1
set /a END_INSTANCE=START_INSTANCE+INSTANCES-1
echo Creando nuevas instancias desde %START_INSTANCE% hasta %END_INSTANCE%.

:: Crear nuevas instancias de frontend
for /L %%i in (%START_INSTANCE%,1,%END_INSTANCE%) do (
    set /a frontendPort=3000+%%i
    echo Creando frontend instance %%i...

    docker run -d --restart unless-stopped --name frontend%%i ^
        -p !frontendPort!:3000 ^
        -v "%cd%\client":/app ^
        -v /app/node_modules ^
        -v /app/.next ^
        --network scrapper-client-network ^
        --ip 10.0.11.1%%i ^
        --cap-add NET_ADMIN ^
        --privileged ^
        --env-file "%cd%\client\.env" ^
        -e NEXT_PUBLIC_API_URL=http://10.0.10.254:8080 ^
        scrapper-client-image:latest
)

echo Todas las nuevas instancias han sido creadas.
pause
