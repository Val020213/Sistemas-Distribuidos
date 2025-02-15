FROM golang:1.23-alpine

WORKDIR /app

EXPOSE ${PORT}

RUN apk add --no-cache iproute2 iptables iputils

COPY server.sh /app

COPY go.mod go.sum ./

RUN go mod download

RUN chmod +x /app/server.sh

CMD ["/bin/sh", "-c", "/app/server.sh && ping 10.0.10.2"]