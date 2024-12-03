FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache git iproute2 iputils

EXPOSE 8080
COPY server.sh /app

RUN chmod +x /app/server.sh

CMD /app/server.sh && ping 10.0.10.2