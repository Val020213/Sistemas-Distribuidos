FROM golang:1.23-alpine

WORKDIR /app

EXPOSE 8080

RUN apk add iproute2 iputils iptables

COPY server.sh /app

RUN chmod +x /app/server.sh

CMD /app/server.sh && ping 10.0.10.2