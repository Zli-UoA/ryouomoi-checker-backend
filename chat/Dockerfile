FROM node:14 as builder
WORKDIR /app
RUN npm install -g typescript
COPY package.json package-lock.json tsconfig.json ./
RUN npm install
COPY . ./
RUN tsc

FROM node:14-alpine as runner
COPY --from=builder /app/lib /app
COPY --from=builder /app/node_modules /app/node_modules
ENTRYPOINT ["node", "/app/index.js"]
