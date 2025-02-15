@echo off

REM 

REM Guardar la imagen my-router-image en un archivo tar
docker save -o scrapper-router-image.tar scrapper-router-image:latest
if %ERRORLEVEL% NEQ 0 (
    echo "Error al guardar la imagen scrapper-router-image"
    exit /b %ERRORLEVEL%
)

REM Guardar la imagen my-backend-image en un archivo tar
docker save -o scrapper-backend-image.tar scrapper-backend-image:latest
if %ERRORLEVEL% NEQ 0 (
    echo "Error al guardar la imagen my-backend-image"
    exit /b %ERRORLEVEL%
)

REM Guardar la imagen my-frontend1-image en un archivo tar
docker save -o scrapper-client-image.tar scrapper-client-image:latest
if %ERRORLEVEL% NEQ 0 (
    echo "Error al guardar la imagen scrapper-client-image"
    exit /b %ERRORLEVEL%
)


echo "Todas las imágenes se han guardado y copiado correctamente."