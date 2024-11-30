# Proyecto de Sistemas Distribuidos

Este es el Proyecto Final de la asignatura Sistemas Distribuidos de la carrera CIencia de la Computación curso 2024-2025. Consiste en una aplicación distribuida para realizar scrapping a páginas web. En el repositorio se incluyen una aplicación de frontend desarrollada con Next.js y un backend desarrollado en Go. La aplicación permite realizar funciones como solicitar scrapping de páginas web para su posterior descarga y brinda también la capacidad de poder descargar páginas que fueron slicitadas por otros usuarios a través de una interfaz web.

## Estructura del proyecto

El proyecto se encuentra dividido en dos carpetas:

* `cliente2`: Contiene el código del cliente frontend en Next.js.
* `go-server`: Contiene el código del backend en Go.

## Requisitos

* Docker
* Docker Compose

## Instalación

1. Clona el repositorio:

   ```bash
   git clone https://github.com/Val020213/Sistemas-Distribuidos
   cd Sistemas-Distribuidos
   ```
2. Dar permisos a los archivo .sh en caso de que nos encontremos en linux

   ```bash
   chmod +x check_directory.sh
   ```
3. Levantar los contenedores

   En el caso de Linux:

   ```bash
   sudo bash run.sh
   ```
   En el caso de WIndows:

   ```cmd
   ./run.bat
   ```
4. Acceder a la aplicación:

   - Cliente 1: [http://localhost:3000](http://localhost:3000)
   - Cliente 2: [http://localhost:3001](http://localhost:3001)
5. Verificar Logs:

   Para ver los logs del frontend, ejecuta:

   ```bash
   docker logs frontend-1
   docker logs frontend-2
   ```
   Para ver los logs del backend, ejecuta:

   ```bash
   dockerlogs backend
   ```
