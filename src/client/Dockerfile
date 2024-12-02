FROM node:20-alpine

WORKDIR /app

COPY package.json ./

COPY package-lock.json ./

ENV NEXT_PUBLIC_API_URL=NEXT_PUBLIC_API_URL

EXPOSE 3000

RUN apk add --no-cache iproute2

COPY client.sh /app

RUN chmod +x /app/client.sh

CMD /app/client.sh

RUN npm install --legacy-peer-deps && npm run dev