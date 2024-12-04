@echo off

cd images || exit /b 1


REM Cargar la imagen my-router-image desde el archivo tar
docker load -i my-router-image.tar
if %ERRORLEVEL% NEQ 0 (
    echo "Error al cargar la imagen my-router-image"
    exit /b %ERRORLEVEL%
)

REM Cargar la imagen my-backend-image desde el archivo tar
docker load -i my-backend-image.tar
if %ERRORLEVEL% NEQ 0 (
    echo "Error al cargar la imagen my-backend-image"
    exit /b %ERRORLEVEL%
)

REM Cargar la imagen my-frontend1-image desde el archivo tar
docker load -i my-frontend1-image.tar
if %ERRORLEVEL% NEQ 0 (
    echo "Error al cargar la imagen my-frontend1-image"
    exit /b %ERRORLEVEL%
)

REM Cargar la imagen my-frontend2-image desde el archivo tar
docker load -i my-frontend2-image.tar
if %ERRORLEVEL% NEQ 0 (
    echo "Error al cargar la imagen my-frontend2-image"
    exit /b %ERRORLEVEL%
)

echo "Todas las imágenes se han cargado correctamente."