FROM node:16-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

# Copy source code
COPY src ./src

EXPOSE 3000

CMD ["npm", "start"] 