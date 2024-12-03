FROM node:20-alpine

WORKDIR /app

COPY package.json ./

COPY package-lock.json ./

ENV NEXT_PUBLIC_API_URL=NEXT_PUBLIC_API_URL

RUN apk add iproute2 iputils && npm install --legacy-peer-deps

EXPOSE 3000

COPY client.sh /app

RUN chmod +x /app/client.sh

CMD /app/client.sh && ping 10.0.10.2