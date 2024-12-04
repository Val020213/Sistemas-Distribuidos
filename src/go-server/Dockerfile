FROM golang:1.23-alpine

WORKDIR /app

EXPOSE 8080

RUN apk add iproute2 iputils iptables

COPY server.sh /app

COPY go.mod go.sum ./

RUN chmod +x /app/server.sh && go mod download

CMD ["/bin/sh", "-c", "/app/server.sh && ping 10.0.10.2"]