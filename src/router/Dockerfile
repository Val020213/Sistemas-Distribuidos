FROM golang:1.23-alpine

WORKDIR /app

RUN echo "net.ipv4.ip_forward=1" | tee -a /etc/sysctl.conf && sysctl -p

# COPY route.sh /root/route.sh

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN chmod +x /app/route.sh

ENTRYPOINT ["/app/route.sh"]