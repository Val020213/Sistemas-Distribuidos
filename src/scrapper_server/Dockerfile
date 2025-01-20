# Usar la imagen base de Go con Alpine
FROM golang:1.23-alpine

# Establecer el directorio de trabajo
WORKDIR /app

# Exponer el puerto en el que el servidor escucha
EXPOSE 8080

# Instalar herramientas necesarias
RUN apk add --no-cache iproute2 iputils iptables

# Copiar el script de inicialización
COPY server.sh /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Asegurarse de que el script de inicio tenga permisos de ejecución
RUN chmod +x /app/server.sh && go mod download

# Comando por defecto para ejecutar el servidor
CMD ["/bin/sh", "-c", "/app/server.sh && ping 10.0.10.2"]
