FROM node:20-alpine

WORKDIR /app

COPY package.json ./

COPY package-lock.json ./

ENV NEXT_PUBLIC_API_URL=NEXT_PUBLIC_API_URL

RUN npm install --legacy-peer-deps

EXPOSE 3000

CMD ["npm", "run", "dev"]