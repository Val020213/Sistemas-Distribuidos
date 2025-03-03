FROM golang:1.23-alpine

WORKDIR /app

RUN echo "net.ipv4.ip_forward=1" | tee -a /etc/sysctl.conf && sysctl -p

EXPOSE 8080

# COPY route.sh /root/route.sh

COPY go.mod go.sum ./

RUN go mod download

# COPY . .

RUN apk add --no-cache iptables && chmod +x /app/route.sh

ENTRYPOINT ["/app/route.sh"]