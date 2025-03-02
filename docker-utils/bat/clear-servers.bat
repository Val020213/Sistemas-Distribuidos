@echo off
REM Check if at least one argument is provided
if "%~1"=="" (
    echo Usage: %0 number1 number2 ...
    exit /b 1
)

for %%i in (%*) do (
    echo Stopping container: backend%%i
    docker stop backend%%i

    echo Stopping container: mongo_bp%%i
    docker stop mongo_bp%%i

    echo Deleting volume: mongo_volume_bp%%i
    docker volume rm mongo_volume_bp%%i
)

echo Pruning stopped containers and unused volumes...
docker container prune -f
docker volume prune -f

echo Done.