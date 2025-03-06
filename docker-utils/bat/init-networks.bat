@echo off

cd ./src || exit 1

rem check scrapper-server-network docker networks existence

docker network inspect scrapper-server-network > nul 2>&1
if %errorlevel% equ 0 (
    echo Network scrapper-server-network exists.
) else (
    docker network create scrapper-server-network --subnet 10.0.10.0/24
    echo Network scrapper-server-network created.
)

rem check scrapper-client-network docker network existence 

docker network inspect scrapper-client-network > nul 2>&1
if %errorlevel% equ 0 (
    echo Network scrapper-client-network exists.
) else (
    docker network create scrapper-client-network --subnet 10.0.11.0/24
    echo Network scrapper-client-network created.
)


rem check router docker image existence 

docker image inspect scrapper-router-image:latest > nul 2>&1
if %errorlevel% equ 0 (
    echo Image scrapper-router-image exists.
) else (
    docker build -t scrapper-router-image:latest ./router
    echo Image scrapper-router-image created.
)

rem check router container existence

docker container inspect router > nul 2>&1
if %errorlevel% equ 0 (
    docker container stop router
    docker container rm router
    echo Container router removed.
)

docker run -d --rm --name router -p 1234:8080 --cap-add NET_ADMIN -v "%cd%\router":/app ^ -v "%cd%\certs":/app/certs -v "go-mod-cache:/go/pkg/mod" -e PYTHONUNBUFFERED=1 scrapper-router-image
echo Container router executed.

docker network connect --ip 10.0.10.254 scrapper-server-network router
docker network connect --ip 10.0.11.254 scrapper-client-network router

@REM docker run -d --rm --name mcproxy --cap-add NET_ADMIN -e PYTHONUNBUFFERED=1 router
@REM echo Container router executed.

@REM docker network connect --ip 10.0.11.253 scrapper-client-network mcproxy
@REM docker network connect --ip 10.0.10.253 scrapper-server-network mcproxy

@REM echo Container router connected to client and server networks

