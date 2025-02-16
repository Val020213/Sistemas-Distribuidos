FROM golang:1.23-alpine

WORKDIR /app

EXPOSE ${PORT}

EXPOSE ${RPC_PORT}

RUN apk add --no-cache iproute2 iptables iputils

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN chmod +x /app/server.sh

CMD ["/bin/sh", "-c", "/app/server.sh"]
