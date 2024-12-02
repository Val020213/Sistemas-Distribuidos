FROM golang:1.23-alpine

RUN apk add --no-cache git

WORKDIR /app

EXPOSE 8080

COPY server.sh /app

RUN chmod +x /app/server.sh

CMD /app/server.sh

RUN ["go", "run", "."]