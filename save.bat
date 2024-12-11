@echo off

REM 
cd images || exit /b 1

REM Guardar la imagen my-router-image en un archivo tar
docker save -o my-router-image.tar my-router-image:latest
if %ERRORLEVEL% NEQ 0 (
    echo "Error al guardar la imagen my-router-image"
    exit /b %ERRORLEVEL%
)

REM Guardar la imagen my-backend-image en un archivo tar
docker save -o my-backend-image.tar my-backend-image:latest
if %ERRORLEVEL% NEQ 0 (
    echo "Error al guardar la imagen my-backend-image"
    exit /b %ERRORLEVEL%
)

REM Guardar la imagen my-frontend1-image en un archivo tar
docker save -o my-frontend1-image.tar my-frontend1-image:latest
if %ERRORLEVEL% NEQ 0 (
    echo "Error al guardar la imagen my-frontend1-image"
    exit /b %ERRORLEVEL%
)

REM Guardar la imagen my-frontend2-image en un archivo tar
docker save -o my-frontend2-image.tar my-frontend2-image:latest
if %ERRORLEVEL% NEQ 0 (
    echo "Error al guardar la imagen my-frontend2-image"
    exit /b %ERRORLEVEL%
)

echo "Todas las imágenes se han guardado y copiado correctamente."