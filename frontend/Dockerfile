# Frontend multi-stage Dockerfile: build with Node, serve with nginx
FROM node:18-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json* ./
# install dependencies
RUN npm ci --silent || npm install --silent
COPY . .
RUN npm run build

# Production image
FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["/bin/sh", "-c", "nginx -g 'daemon off;'"]
