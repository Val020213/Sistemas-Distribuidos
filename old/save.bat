@REM @echo off

@REM REM 

@REM REM Guardar la imagen my-router-image en un archivo tar
@REM docker save -o scrapper-router-image.tar scrapper-router-image:latest
@REM if %ERRORLEVEL% NEQ 0 (
@REM     echo "Error al guardar la imagen scrapper-router-image"
@REM     exit /b %ERRORLEVEL%
@REM )

@REM REM Guardar la imagen my-backend-image en un archivo tar
@REM docker save -o scrapper-backend-image.tar scrapper-backend-image:latest
@REM if %ERRORLEVEL% NEQ 0 (
@REM     echo "Error al guardar la imagen my-backend-image"
@REM     exit /b %ERRORLEVEL%
@REM )

@REM REM Guardar la imagen my-frontend1-image en un archivo tar
@REM docker save -o scrapper-client-image.tar scrapper-client-image:latest
@REM if %ERRORLEVEL% NEQ 0 (
@REM     echo "Error al guardar la imagen scrapper-client-image"
@REM     exit /b %ERRORLEVEL%
@REM )


@REM echo "Todas las imágenes se han guardado y copiado correctamente."

docker sa