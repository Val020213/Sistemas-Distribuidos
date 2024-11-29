# TODO : MAKE THE README FILE

# TODO : Is necessary make the entrypoint.sh file?

# TODO : Why in Linux the volume work and in Windows not?

# TODO : Testing the container

# TODO : The run files is not tested yet

---

# THIS IS A TEMPLATE FOR A README FILE MADE BY GTP4o

Para crear un informe del proyecto en el archivo

README.md

, puedes seguir una estructura que incluya una descripción general del proyecto, instrucciones de instalación, uso, estructura del proyecto, y cualquier otra información relevante. Aquí tienes un ejemplo de cómo podrías estructurar tu

README.md

:

````markdown
# Proyecto de Sistemas Distribuidos

## Descripción

Este proyecto consiste en una aplicación distribuida que incluye un frontend desarrollado en Next.js y un backend desarrollado en Go. La aplicación permite realizar diversas operaciones, como la obtención de datos y la descarga de archivos, a través de una interfaz web.

## Estructura del Proyecto

El proyecto está organizado en las siguientes carpetas:

- `cliente1`: Contiene el código del primer cliente frontend en Next.js.
- `cliente2`: Contiene el código del segundo cliente frontend en Next.js.
- `go-server`: Contiene el código del backend en Go.

## Requisitos

- Docker
- Docker Compose

## Instalación

1. Clona el repositorio:

   ```sh
   git clone https://github.com/tu-usuario/tu-repositorio.git
   cd tu-repositorio
   ```
````

2. Crea un archivo `.env` en la raíz del proyecto con las siguientes variables de entorno:

   ```env
   SERVER_HOST=0.0.0.0
   SERVER_PORT=8080
   ```

3. Crea un archivo `.env.local` en las carpetas `cliente1` y `cliente2` con la siguiente variable de entorno:

   ```env
   NEXT_PUBLIC_API_URL=http://backend:8080
   ```

## Uso

### Levantar los contenedores

Para levantar los contenedores, ejecuta el siguiente comando:

```sh
docker-compose up --build
```

### Acceder a la aplicación

- Cliente 1: [http://localhost:3000](http://localhost:3000)
- Cliente 2: [http://localhost:3001](http://localhost:3001)

### Verificar los logs

Para ver los logs del backend, ejecuta:

```sh
docker logs <backend_container_id_or_name>
```

Para ver los logs de los frontends, ejecuta:

```sh
docker logs <frontend1_container_id_or_name>
docker logs <frontend2_container_id_or_name>
```

## Scripts

### Script Bash (`run.sh`)

```sh
#!/bin/bash

# Levantar los contenedores con Docker Compose
docker-compose up --build
```

### Script Batch (`run.bat`)

```bat
@echo off
REM Levantar los contenedores con Docker Compose
docker-compose up --build
```

### Script Bash para abrir el navegador (`open_browser.sh`)

```sh
#!/bin/bash

# Esperar a que los contenedores estén listos
sleep 10

# Abrir el navegador para el frontend1
xdg-open http://localhost:3000

# Abrir el navegador para el frontend2
xdg-open http://localhost:3001
```

### Script Batch para abrir el navegador (`open_browser.bat`)

```bat
@echo off
REM Esperar a que los contenedores estén listos
timeout /t 10

REM Abrir el navegador para el frontend1
start http://localhost:3000

REM Abrir el navegador para el frontend2
start http://localhost:3001
```

## Contribuciones

Las contribuciones son bienvenidas. Por favor, abre un issue o un pull request para discutir cualquier cambio que desees realizar.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT. Consulta el archivo [LICENSE](LICENSE) para obtener más detalles.

```

Este

README.md

 proporciona una descripción completa del proyecto, instrucciones de instalación y uso, y detalles sobre la estructura del proyecto y los scripts disponibles. Puedes personalizarlo según las necesidades específicas de tu proyecto.
```
