@echo off

REM Navegar al directorio que contiene el docker-compose.yml
cd src || exit /b 1

REM 
docker-compose down -v 

REM 
docker-compose up --build


REM Dar tiempo para que los contenedores se levanten
timeout /t 10

REM Abrir el navegador para el frontend1
start http://localhost:3000

REM Abrir el navegador para el frontend2
start http://localhost:3001